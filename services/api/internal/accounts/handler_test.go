package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/montra/montra/services/api/internal/auth"
)

type mockAccountRepo struct {
	accounts map[string]*Account
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: make(map[string]*Account)}
}

func (m *mockAccountRepo) ListWithBalance(ctx context.Context, userID string) ([]Account, error) {
	var list []Account
	for _, a := range m.accounts {
		if a.UserID == userID && a.ArchivedAt == nil {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockAccountRepo) GetByIDWithBalance(ctx context.Context, userID, accountID string) (*Account, error) {
	a, ok := m.accounts[accountID]
	if !ok || a.UserID != userID {
		return nil, ErrAccountNotFound
	}
	return a, nil
}

func (m *mockAccountRepo) Create(ctx context.Context, a *Account) error {
	a.ID = "acc-uuid-1"
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
	a.CurrentBalance = a.OpeningBalance
	m.accounts[a.ID] = a
	return nil
}

func (m *mockAccountRepo) Update(ctx context.Context, userID, accountID, name, accountType string) error {
	a, ok := m.accounts[accountID]
	if !ok || a.UserID != userID || a.ArchivedAt != nil {
		return ErrAccountNotFound
	}
	if name != "" {
		a.Name = name
	}
	if accountType != "" {
		a.AccountType = accountType
	}
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *mockAccountRepo) Archive(ctx context.Context, userID, accountID string) error {
	a, ok := m.accounts[accountID]
	if !ok || a.UserID != userID || a.ArchivedAt != nil {
		return ErrAccountNotFound
	}
	now := time.Now().UTC()
	a.ArchivedAt = &now
	return nil
}

func setupTestAccountRouter() (http.Handler, *mockAccountRepo, *auth.TokenManager) {
	repo := newMockAccountRepo()
	tm, _ := auth.NewTokenManager("test-jwt-secret-at-least-32-characters-long", 2*time.Hour)
	svc := NewService(repo)
	handler := NewHandler(svc, tm)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux, repo, tm
}

func TestAccountEndpoints(t *testing.T) {
	mux, repo, tm := setupTestAccountRouter()

	userAToken, _ := tm.Generate("usr-A", "userA@montra.id")
	userBToken, _ := tm.Generate("usr-B", "userB@montra.id")

	// 1. Create account for User A
	createPayload := map[string]any{
		"name":            "Dompet Tunai",
		"account_type":    "cash",
		"opening_balance": int64(250000),
	}
	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/v1/accounts", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+userAToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on create account, got %d: %s", w.Code, w.Body.String())
	}

	var createResp struct {
		Data AccountDTO `json:"data"`
	}
	_ = json.NewDecoder(w.Body).Decode(&createResp)
	if createResp.Data.CurrentBalance != 250000 {
		t.Errorf("expected balance 250000, got %d", createResp.Data.CurrentBalance)
	}
	accID := createResp.Data.ID

	// 2. Reject negative opening balance
	badPayload := map[string]any{
		"name":            "Hutang",
		"account_type":    "cash",
		"opening_balance": int64(-50000),
	}
	badBody, _ := json.Marshal(badPayload)
	reqBad := httptest.NewRequest("POST", "/api/v1/accounts", bytes.NewReader(badBody))
	reqBad.Header.Set("Authorization", "Bearer "+userAToken)
	wBad := httptest.NewRecorder()
	mux.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for negative opening balance, got %d", wBad.Code)
	}

	// 3. User A can get their account
	reqGet := httptest.NewRequest("GET", "/api/v1/accounts/"+accID, nil)
	reqGet.Header.Set("Authorization", "Bearer "+userAToken)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for account owner, got %d", wGet.Code)
	}

	// 4. IDOR Protection: User B CANNOT get User A's account
	reqIDOR := httptest.NewRequest("GET", "/api/v1/accounts/"+accID, nil)
	reqIDOR.Header.Set("Authorization", "Bearer "+userBToken)
	wIDOR := httptest.NewRecorder()
	mux.ServeHTTP(wIDOR, reqIDOR)
	if wIDOR.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found for non-owner User B (IDOR prevention), got %d", wIDOR.Code)
	}

	// 5. Update account
	updatePayload := map[string]any{
		"name": "Dompet Utama",
	}
	updBody, _ := json.Marshal(updatePayload)
	reqUpd := httptest.NewRequest("PATCH", "/api/v1/accounts/"+accID, bytes.NewReader(updBody))
	reqUpd.Header.Set("Authorization", "Bearer "+userAToken)
	wUpd := httptest.NewRecorder()
	mux.ServeHTTP(wUpd, reqUpd)
	if wUpd.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on update, got %d", wUpd.Code)
	}

	// 6. Archive account
	reqDel := httptest.NewRequest("DELETE", "/api/v1/accounts/"+accID, nil)
	reqDel.Header.Set("Authorization", "Bearer "+userAToken)
	wDel := httptest.NewRecorder()
	mux.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on archive, got %d", wDel.Code)
	}

	// 7. Archived account should not appear in active list
	reqList := httptest.NewRequest("GET", "/api/v1/accounts", nil)
	reqList.Header.Set("Authorization", "Bearer "+userAToken)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)

	var listResp struct {
		Data []AccountDTO `json:"data"`
	}
	_ = json.NewDecoder(wList.Body).Decode(&listResp)
	if len(listResp.Data) != 0 {
		t.Errorf("expected 0 active accounts after archive, got %d", len(listResp.Data))
	}
	_ = repo
}
