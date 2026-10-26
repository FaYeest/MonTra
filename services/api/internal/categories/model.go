package categories

import (
	"errors"
	"time"
)

var (
	// ErrCategoryNotFound indicates category was not found or is not accessible to the user.
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
	// ErrInvalidCategoryName indicates invalid name length.
	ErrInvalidCategoryName = errors.New("nama kategori harus antara 1 sampai 50 karakter")
	// ErrInvalidTransactionType indicates transaction_type is not 'income' or 'expense'.
	ErrInvalidTransactionType = errors.New("tipe transaksi harus 'income' atau 'expense'")
)

// Category represents a financial category entity.
type Category struct {
	ID              string     `json:"id"`
	UserID          *string    `json:"user_id,omitempty"` // NULL for system categories
	Name            string     `json:"name"`
	TransactionType string     `json:"transaction_type"` // 'income' or 'expense'
	Icon            string     `json:"icon"`
	Color           string     `json:"color"`
	IsSystem        bool       `json:"is_system"`
	IsArchived      bool       `json:"is_archived"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CategoryDTO is the client-facing representation of a category.
type CategoryDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	TransactionType string `json:"transaction_type"`
	Icon            string `json:"icon"`
	Color           string `json:"color"`
	IsSystem        bool   `json:"is_system"`
}

// ToDTO converts an internal Category model to CategoryDTO.
func (c *Category) ToDTO() CategoryDTO {
	return CategoryDTO{
		ID:              c.ID,
		Name:            c.Name,
		TransactionType: c.TransactionType,
		Icon:            c.Icon,
		Color:           c.Color,
		IsSystem:        c.IsSystem,
	}
}
