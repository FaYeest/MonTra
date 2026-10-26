package transactions

import (
	"context"
	"strings"
	"time"
)

// CreateTransactionRequest represents input for creating a transaction.
type CreateTransactionRequest struct {
	AccountID            string  `json:"account_id"`
	CategoryID           *string `json:"category_id,omitempty"`
	Type                 string  `json:"type"`
	Amount               int64   `json:"amount"`
	Currency             string  `json:"currency,omitempty"`
	Description          string  `json:"description"`
	TransactionDate      string  `json:"transaction_date,omitempty"`
	Notes                *string `json:"notes,omitempty"`
	Source               string  `json:"source,omitempty"`
	CategorySource       string  `json:"category_source,omitempty"`
	CategoryModelVersion *string `json:"category_model_version,omitempty"`
	IdempotencyKey       *string `json:"idempotency_key,omitempty"`
}

// UpdateTransactionRequest represents input for modifying a transaction.
type UpdateTransactionRequest struct {
	AccountID       string  `json:"account_id"`
	CategoryID      *string `json:"category_id,omitempty"`
	Type            string  `json:"type"`
	Amount          int64   `json:"amount"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	Notes           *string `json:"notes,omitempty"`
	CategorySource  string  `json:"category_source,omitempty"`
}

// Service provides business logic for transactions.
type Service struct {
	repo Repository
}

// NewService constructs a Transaction Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateTransaction validates inputs and creates a new transaction, respecting idempotency.
func (s *Service) CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*TransactionDetailDTO, error) {
	// 1. Idempotency check
	if req.IdempotencyKey != nil && strings.TrimSpace(*req.IdempotencyKey) != "" {
		key := strings.TrimSpace(*req.IdempotencyKey)
		existing, err := s.repo.GetByIdempotencyKey(ctx, userID, key)
		if err == nil && existing != nil {
			return existing, nil
		}
	}

	// 2. Amount validation
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// 3. Type validation
	txType := strings.ToLower(strings.TrimSpace(req.Type))
	if txType != "income" && txType != "expense" {
		return nil, ErrInvalidType
	}

	// 4. Description validation
	desc := strings.TrimSpace(req.Description)
	if len(desc) < 1 || len(desc) > 120 {
		return nil, ErrInvalidDescription
	}

	// 5. Account ownership validation
	accountID := strings.TrimSpace(req.AccountID)
	if err := s.repo.ValidateAccount(ctx, userID, accountID); err != nil {
		return nil, err
	}

	// 6. Category validation
	var categoryID *string
	if req.CategoryID != nil && strings.TrimSpace(*req.CategoryID) != "" {
		catID := strings.TrimSpace(*req.CategoryID)
		if err := s.repo.ValidateCategory(ctx, userID, catID, txType); err != nil {
			return nil, err
		}
		categoryID = &catID
	}

	// 7. Date parsing
	txDate := time.Now().UTC()
	if strings.TrimSpace(req.TransactionDate) != "" {
		parsedDate, err := time.Parse(time.RFC3339, strings.TrimSpace(req.TransactionDate))
		if err != nil {
			// fallback YYYY-MM-DD
			parsedDate, err = time.Parse("2006-01-02", strings.TrimSpace(req.TransactionDate))
			if err != nil {
				return nil, err
			}
		}
		txDate = parsedDate.UTC()
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "IDR"
	}

	source := strings.ToLower(strings.TrimSpace(req.Source))
	if source == "" {
		source = "manual"
	}

	catSource := strings.ToLower(strings.TrimSpace(req.CategorySource))
	if catSource == "" {
		catSource = "manual"
	}

	var idempotencyKey *string
	if req.IdempotencyKey != nil && strings.TrimSpace(*req.IdempotencyKey) != "" {
		key := strings.TrimSpace(*req.IdempotencyKey)
		idempotencyKey = &key
	}

	t := &Transaction{
		UserID:               userID,
		AccountID:            accountID,
		CategoryID:           categoryID,
		Type:                 txType,
		Amount:               req.Amount,
		Currency:             currency,
		Description:          desc,
		TransactionDate:      txDate,
		Notes:                req.Notes,
		Source:               source,
		CategorySource:       catSource,
		CategoryModelVersion: req.CategoryModelVersion,
		IdempotencyKey:       idempotencyKey,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, userID, t.ID)
}

// GetTransaction retrieves a transaction by ID.
func (s *Service) GetTransaction(ctx context.Context, userID, transactionID string) (*TransactionDetailDTO, error) {
	return s.repo.GetByID(ctx, userID, transactionID)
}

// ListTransactions retrieves filtered transactions.
func (s *Service) ListTransactions(ctx context.Context, userID string, f Filter) ([]TransactionDetailDTO, int, error) {
	return s.repo.List(ctx, userID, f)
}

// UpdateTransaction updates an existing transaction.
func (s *Service) UpdateTransaction(ctx context.Context, userID, transactionID string, req UpdateTransactionRequest) (*TransactionDetailDTO, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	txType := strings.ToLower(strings.TrimSpace(req.Type))
	if txType != "income" && txType != "expense" {
		return nil, ErrInvalidType
	}

	desc := strings.TrimSpace(req.Description)
	if len(desc) < 1 || len(desc) > 120 {
		return nil, ErrInvalidDescription
	}

	accountID := strings.TrimSpace(req.AccountID)
	if err := s.repo.ValidateAccount(ctx, userID, accountID); err != nil {
		return nil, err
	}

	var categoryID *string
	if req.CategoryID != nil && strings.TrimSpace(*req.CategoryID) != "" {
		catID := strings.TrimSpace(*req.CategoryID)
		if err := s.repo.ValidateCategory(ctx, userID, catID, txType); err != nil {
			return nil, err
		}
		categoryID = &catID
	}

	parsedDate, err := time.Parse(time.RFC3339, strings.TrimSpace(req.TransactionDate))
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02", strings.TrimSpace(req.TransactionDate))
		if err != nil {
			return nil, err
		}
	}

	catSource := strings.ToLower(strings.TrimSpace(req.CategorySource))
	if catSource == "" {
		catSource = "manual"
	}

	t := &Transaction{
		AccountID:       accountID,
		CategoryID:      categoryID,
		Type:            txType,
		Amount:          req.Amount,
		Description:     desc,
		TransactionDate: parsedDate.UTC(),
		Notes:           req.Notes,
		CategorySource:  catSource,
	}

	if err := s.repo.Update(ctx, userID, transactionID, t); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, userID, transactionID)
}

// DeleteTransaction removes a transaction.
func (s *Service) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	return s.repo.Delete(ctx, userID, transactionID)
}
