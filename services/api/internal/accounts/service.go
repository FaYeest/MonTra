package accounts

import (
	"context"
	"errors"
	"strings"
)

// CreateAccountRequest holds payload to create an account.
type CreateAccountRequest struct {
	Name           string `json:"name"`
	AccountType    string `json:"account_type"`
	OpeningBalance int64  `json:"opening_balance"`
	Currency       string `json:"currency,omitempty"`
}

// UpdateAccountRequest holds payload to edit an account.
type UpdateAccountRequest struct {
	Name        string `json:"name,omitempty"`
	AccountType string `json:"account_type,omitempty"`
}

// Service handles business logic and validation for accounts.
type Service struct {
	repo Repository
}

// NewService constructs an Account Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListAccounts returns all active accounts for a user with running balances.
func (s *Service) ListAccounts(ctx context.Context, userID string) ([]AccountDTO, error) {
	accs, err := s.repo.ListWithBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]AccountDTO, len(accs))
	for i, a := range accs {
		result[i] = a.ToDTO()
	}
	return result, nil
}

// GetAccount retrieves a single account by ID with running balance.
func (s *Service) GetAccount(ctx context.Context, userID, accountID string) (*AccountDTO, error) {
	acc, err := s.repo.GetByIDWithBalance(ctx, userID, accountID)
	if err != nil {
		return nil, err
	}
	dto := acc.ToDTO()
	return &dto, nil
}

// CreateAccount validates input and persists a new account.
func (s *Service) CreateAccount(ctx context.Context, userID string, req CreateAccountRequest) (*AccountDTO, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) < 1 || len(name) > 50 {
		return nil, ErrInvalidAccountName
	}

	accType := strings.ToLower(strings.TrimSpace(req.AccountType))
	if !ValidAccountTypes[accType] {
		return nil, ErrInvalidAccountType
	}

	if req.OpeningBalance < 0 {
		return nil, ErrNegativeOpeningBalance
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "IDR"
	}
	if len(currency) != 3 {
		return nil, errors.New("currency harus 3 karakter ISO 4217")
	}

	acc := &Account{
		UserID:         userID,
		Name:           name,
		AccountType:    accType,
		OpeningBalance: req.OpeningBalance,
		Currency:       currency,
	}

	if err := s.repo.Create(ctx, acc); err != nil {
		return nil, err
	}

	dto := acc.ToDTO()
	return &dto, nil
}

// UpdateAccount validates changes and updates the account.
func (s *Service) UpdateAccount(ctx context.Context, userID, accountID string, req UpdateAccountRequest) (*AccountDTO, error) {
	name := strings.TrimSpace(req.Name)
	if name != "" && (len(name) < 1 || len(name) > 50) {
		return nil, ErrInvalidAccountName
	}

	accType := strings.ToLower(strings.TrimSpace(req.AccountType))
	if accType != "" && !ValidAccountTypes[accType] {
		return nil, ErrInvalidAccountType
	}

	if err := s.repo.Update(ctx, userID, accountID, name, accType); err != nil {
		return nil, err
	}

	return s.GetAccount(ctx, userID, accountID)
}

// ArchiveAccount archives an account so it no longer appears in active lists.
func (s *Service) ArchiveAccount(ctx context.Context, userID, accountID string) error {
	return s.repo.Archive(ctx, userID, accountID)
}
