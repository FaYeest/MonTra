package transfers

import (
	"context"
	"strings"
	"time"
)

// CreateTransferRequest holds parameters for recording an internal fund transfer.
type CreateTransferRequest struct {
	SourceAccountID      string  `json:"source_account_id"`
	DestinationAccountID string  `json:"destination_account_id"`
	Amount               int64   `json:"amount"`
	Currency             string  `json:"currency,omitempty"`
	FeeAmount            int64   `json:"fee_amount"`
	TransactionDate      string  `json:"transaction_date,omitempty"`
	Description          *string `json:"description,omitempty"`
	IdempotencyKey       *string `json:"idempotency_key,omitempty"`
}

// Service provides business rules and orchestrates internal transfers.
type Service struct {
	repo Repository
}

// NewService constructs a Transfer Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateTransfer executes an internal fund transfer atomically.
func (s *Service) CreateTransfer(ctx context.Context, userID string, req CreateTransferRequest) (*TransferDetailDTO, error) {
	// 1. Idempotency check
	if req.IdempotencyKey != nil && strings.TrimSpace(*req.IdempotencyKey) != "" {
		key := strings.TrimSpace(*req.IdempotencyKey)
		existing, err := s.repo.GetByIdempotencyKey(ctx, userID, key)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	sourceID := strings.TrimSpace(req.SourceAccountID)
	destID := strings.TrimSpace(req.DestinationAccountID)

	if sourceID == "" || destID == "" {
		return nil, ErrAccountNotFound
	}

	if sourceID == destID {
		return nil, ErrSameAccountTransfer
	}

	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if req.FeeAmount < 0 {
		return nil, ErrInvalidFee
	}

	txDate, err := parseTransferDate(req.TransactionDate)
	if err != nil {
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "IDR"
	}

	var idempotencyKey *string
	if req.IdempotencyKey != nil && strings.TrimSpace(*req.IdempotencyKey) != "" {
		key := strings.TrimSpace(*req.IdempotencyKey)
		idempotencyKey = &key
	}

	t := &Transfer{
		UserID:               userID,
		SourceAccountID:      sourceID,
		DestinationAccountID: destID,
		Amount:               req.Amount,
		Currency:             currency,
		FeeAmount:            req.FeeAmount,
		TransactionDate:      txDate,
		Description:          req.Description,
		IdempotencyKey:       idempotencyKey,
	}

	if err := s.repo.CreateAtomic(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, userID, t.ID)
}

// GetTransfer retrieves a single transfer by ID.
func (s *Service) GetTransfer(ctx context.Context, userID, transferID string) (*TransferDetailDTO, error) {
	return s.repo.GetByID(ctx, userID, transferID)
}

// ListTransfers returns paginated transfers for the user.
func (s *Service) ListTransfers(ctx context.Context, userID string, limit, offset int) ([]TransferDetailDTO, int, error) {
	return s.repo.List(ctx, userID, limit, offset)
}

// DeleteTransfer deletes a transfer record.
func (s *Service) DeleteTransfer(ctx context.Context, userID, transferID string) error {
	return s.repo.Delete(ctx, userID, transferID)
}

func parseTransferDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Now().UTC(), nil
	}
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05.999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, raw); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, ErrInvalidDate
}
