package transfers

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

type mockTransferRepo struct {
	transfers map[string]*TransferDetailDTO
	accounts  map[string]string // accountID -> userID
}

func newMockTransferRepo() *mockTransferRepo {
	return &mockTransferRepo{
		transfers: make(map[string]*TransferDetailDTO),
		accounts: map[string]string{
			"acc-cash": "usr-1",
			"acc-bank": "usr-1",
		},
	}
}

func (m *mockTransferRepo) CreateAtomic(ctx context.Context, t *Transfer) error {
	u1, ok1 := m.accounts[t.SourceAccountID]
	u2, ok2 := m.accounts[t.DestinationAccountID]
	if !ok1 || !ok2 || u1 != t.UserID || u2 != t.UserID {
		return ErrAccountNotFound
	}
	t.ID = "tr-uuid-123"
	dto := &TransferDetailDTO{
		ID:                     t.ID,
		SourceAccountID:        t.SourceAccountID,
		SourceAccountName:      "Dompet Tunai",
		DestinationAccountID:   t.DestinationAccountID,
		DestinationAccountName: "BCA Tabungan",
		Amount:                 t.Amount,
		Currency:               t.Currency,
		FeeAmount:              t.FeeAmount,
		TransactionDate:        t.TransactionDate.UTC().Format(time.RFC3339),
		Description:            t.Description,
		CreatedAt:              time.Now().UTC().Format(time.RFC3339),
	}
	m.transfers[t.ID] = dto
	return nil
}

func (m *mockTransferRepo) GetByID(ctx context.Context, userID, transferID string) (*TransferDetailDTO, error) {
	item, ok := m.transfers[transferID]
	if !ok {
		return nil, ErrTransferNotFound
	}
	return item, nil
}

func (m *mockTransferRepo) GetByIdempotencyKey(ctx context.Context, userID, key string) (*TransferDetailDTO, error) {
	return nil, ErrTransferNotFound
}

func (m *mockTransferRepo) List(ctx context.Context, userID string, limit, offset int) ([]TransferDetailDTO, int, error) {
	var list []TransferDetailDTO
	for _, tr := range m.transfers {
		list = append(list, *tr)
	}
	return list, len(list), nil
}

func (m *mockTransferRepo) Delete(ctx context.Context, userID, transferID string) error {
	if _, ok := m.transfers[transferID]; !ok {
		return ErrTransferNotFound
	}
	delete(m.transfers, transferID)
	return nil
}

func TestTransferEndpoints(t *testing.T) {
	repo := newMockTransferRepo()
	tm, _ := auth.NewTokenManager("test-jwt-secret-at-least-32-characters-long", 2*time.Hour)
	svc := NewService(repo)
	handler := NewHandler(svc, tm)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	userToken, _ := tm.Generate("usr-1", "user@montra.id")

	// 1. Successful transfer
	desc := "Top-up dari kas"
	createReq := CreateTransferRequest{
		SourceAccountID:      "acc-cash",
		DestinationAccountID: "acc-bank",
		Amount:               100000,
		FeeAmount:            2500,
		Description:          &desc,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/transfers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on transfer, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Reject same source and destination account
	sameAccountReq := createReq
	sameAccountReq.DestinationAccountID = "acc-cash"
	saJSON, _ := json.Marshal(sameAccountReq)
	reqSA := httptest.NewRequest("POST", "/api/v1/transfers", bytes.NewReader(saJSON))
	reqSA.Header.Set("Authorization", "Bearer "+userToken)
	wSA := httptest.NewRecorder()
	mux.ServeHTTP(wSA, reqSA)
	if wSA.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for same account transfer, got %d", wSA.Code)
	}

	// 3. Reject negative fee
	badFeeReq := createReq
	badFeeReq.FeeAmount = -500
	bfJSON, _ := json.Marshal(badFeeReq)
	reqBF := httptest.NewRequest("POST", "/api/v1/transfers", bytes.NewReader(bfJSON))
	reqBF.Header.Set("Authorization", "Bearer "+userToken)
	wBF := httptest.NewRecorder()
	mux.ServeHTTP(wBF, reqBF)
	if wBF.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for negative fee, got %d", wBF.Code)
	}

	// 4. List transfers
	reqList := httptest.NewRequest("GET", "/api/v1/transfers", nil)
	reqList.Header.Set("Authorization", "Bearer "+userToken)
	wList := httptest.NewRecorder()
	mux.ServeHTTP(wList, reqList)

	var listResp struct {
		Data []TransferDetailDTO `json:"data"`
	}
	_ = json.NewDecoder(wList.Body).Decode(&listResp)
	if len(listResp.Data) != 1 {
		t.Errorf("expected 1 transfer record, got %d", len(listResp.Data))
	}
}
