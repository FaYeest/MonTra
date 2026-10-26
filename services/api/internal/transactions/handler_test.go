package transactions

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

type mockTxRepo struct {
	transactions map[string]*TransactionDetailDTO
	accounts     map[string]string // accountID -> userID
	categories   map[string]string // categoryID -> type
}

func newMockTxRepo() *mockTxRepo {
	return &mockTxRepo{
		transactions: make(map[string]*TransactionDetailDTO),
		accounts:     map[string]string{"acc-1": "usr-1"},
		categories:   map[string]string{"cat-food": "expense"},
	}
}

func (m *mockTxRepo) ValidateAccount(ctx context.Context, userID, accountID string) error {
	owner, ok := m.accounts[accountID]
	if !ok || owner != userID {
		return ErrAccountNotFound
	}
	return nil
}

func (m *mockTxRepo) ValidateCategory(ctx context.Context, userID, categoryID, txType string) error {
	catType, ok := m.categories[categoryID]
	if !ok || catType != txType {
		return ErrCategoryNotValid
	}
	return nil
}

func (m *mockTxRepo) Create(ctx context.Context, t *Transaction) error {
	t.ID = "tx-123"
	dto := &TransactionDetailDTO{
		ID:              t.ID,
		AccountID:       t.AccountID,
		AccountName:     "Dompet Tunai",
		CategoryID:      t.CategoryID,
		Type:            t.Type,
		Amount:          t.Amount,
		Currency:        t.Currency,
		Description:     t.Description,
		TransactionDate: t.TransactionDate.UTC().Format(time.RFC3339),
		Source:          t.Source,
		CategorySource:  t.CategorySource,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
	m.transactions[t.ID] = dto
	return nil
}

func (m *mockTxRepo) GetByID(ctx context.Context, userID, transactionID string) (*TransactionDetailDTO, error) {
	item, ok := m.transactions[transactionID]
	if !ok {
		return nil, ErrTransactionNotFound
	}
	return item, nil
}

func (m *mockTxRepo) GetByIdempotencyKey(ctx context.Context, userID, key string) (*TransactionDetailDTO, error) {
	for _, item := range m.transactions {
		if item.Description == key { // simple mock lookup
			return item, nil
		}
	}
	return nil, ErrTransactionNotFound
}

func (m *mockTxRepo) List(ctx context.Context, userID string, f Filter) ([]TransactionDetailDTO, int, error) {
	var list []TransactionDetailDTO
	for _, item := range m.transactions {
		if f.Type == "" || item.Type == f.Type {
			list = append(list, *item)
		}
	}
	return list, len(list), nil
}

func (m *mockTxRepo) Update(ctx context.Context, userID, transactionID string, t *Transaction) error {
	item, ok := m.transactions[transactionID]
	if !ok {
		return ErrTransactionNotFound
	}
	item.Amount = t.Amount
	item.Description = t.Description
	return nil
}

func (m *mockTxRepo) Delete(ctx context.Context, userID, transactionID string) error {
	if _, ok := m.transactions[transactionID]; !ok {
		return ErrTransactionNotFound
	}
	delete(m.transactions, transactionID)
	return nil
}

func TestTransactionEndpoints(t *testing.T) {
	repo := newMockTxRepo()
	tm, _ := auth.NewTokenManager("test-jwt-secret-at-least-32-characters-long", 2*time.Hour)
	svc := NewService(repo)
	handler := NewHandler(svc, tm)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	userToken, _ := tm.Generate("usr-1", "user@montra.id")
	otherToken, _ := tm.Generate("usr-2", "other@montra.id")

	// 1. Create valid expense transaction
	catID := "cat-food"
	createBody := CreateTransactionRequest{
		AccountID:       "acc-1",
		CategoryID:      &catID,
		Type:            "expense",
		Amount:          25000,
		Description:     "Kopi Kenangan",
		TransactionDate: "2026-10-26T10:00:00Z",
	}
	body, _ := json.Marshal(createBody)
	req := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on create transaction, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Reject non-positive amount
	badAmountBody := createBody
	badAmountBody.Amount = 0
	bJSON, _ := json.Marshal(badAmountBody)
	reqBad := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(bJSON))
	reqBad.Header.Set("Authorization", "Bearer "+userToken)
	wBad := httptest.NewRecorder()
	mux.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for amount <= 0, got %d", wBad.Code)
	}

	// 3. Reject account not owned by user (usr-2 tries to use acc-1)
	reqOther := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(body))
	reqOther.Header.Set("Authorization", "Bearer "+otherToken)
	wOther := httptest.NewRecorder()
	mux.ServeHTTP(wOther, reqOther)
	if wOther.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unowned account, got %d", wOther.Code)
	}

	// 4. List transactions
	reqList := httptest.NewRequest("GET", "/api/v1/transactions?type=expense", nil)
	reqList.Header.Set("Authorization", "Bearer "+userToken)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)

	var listResp struct {
		Data []TransactionDetailDTO `json:"data"`
		Meta map[string]any         `json:"meta"`
	}
	_ = json.NewDecoder(wList.Body).Decode(&listResp)
	if len(listResp.Data) != 1 {
		t.Errorf("expected 1 expense transaction, got %d", len(listResp.Data))
	}

	// 5. Delete transaction
	reqDel := httptest.NewRequest("DELETE", "/api/v1/transactions/tx-123", nil)
	reqDel.Header.Set("Authorization", "Bearer "+userToken)
	wDel := httptest.NewRecorder()
	mux.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusOK {
		t.Errorf("expected 200 on delete, got %d", wDel.Code)
	}
}
