package categories

import (
	"context"
	"strings"
)

// CreateCategoryRequest holds input data for adding a user custom category.
type CreateCategoryRequest struct {
	Name            string `json:"name"`
	TransactionType string `json:"transaction_type"`
	Icon            string `json:"icon,omitempty"`
	Color           string `json:"color,omitempty"`
}

// Service provides business logic for category management.
type Service struct {
	repo Repository
}

// NewService constructs a Category service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListCategories returns active categories for the given user, optionally filtered by 'income' or 'expense'.
func (s *Service) ListCategories(ctx context.Context, userID string, txType string) ([]CategoryDTO, error) {
	txType = strings.ToLower(strings.TrimSpace(txType))
	if txType != "" && txType != "income" && txType != "expense" {
		return nil, ErrInvalidTransactionType
	}

	cats, err := s.repo.List(ctx, userID, txType)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryDTO, len(cats))
	for i, c := range cats {
		result[i] = c.ToDTO()
	}
	return result, nil
}

// CreateCategory validates and saves a new user-defined category.
func (s *Service) CreateCategory(ctx context.Context, userID string, req CreateCategoryRequest) (*CategoryDTO, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) < 1 || len(name) > 50 {
		return nil, ErrInvalidCategoryName
	}

	txType := strings.ToLower(strings.TrimSpace(req.TransactionType))
	if txType != "income" && txType != "expense" {
		return nil, ErrInvalidTransactionType
	}

	c := &Category{
		UserID:          &userID,
		Name:            name,
		TransactionType: txType,
		Icon:            strings.TrimSpace(req.Icon),
		Color:           strings.TrimSpace(req.Color),
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	dto := c.ToDTO()
	return &dto, nil
}
