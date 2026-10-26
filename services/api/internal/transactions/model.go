package transactions

import (
	"errors"
	"time"
)

var (
	// ErrTransactionNotFound indicates the transaction was not found or is not owned by the user.
	ErrTransactionNotFound = errors.New("transaksi tidak ditemukan")
	// ErrInvalidAmount indicates non-positive transaction amount.
	ErrInvalidAmount = errors.New("nominal transaksi harus lebih besar dari 0")
	// ErrInvalidType indicates transaction type is not 'income' or 'expense'.
	ErrInvalidType = errors.New("tipe transaksi harus 'income' atau 'expense'")
	// ErrInvalidDescription indicates empty or oversized description.
	ErrInvalidDescription = errors.New("deskripsi transaksi harus antara 1 sampai 120 karakter")
	// ErrAccountNotFound indicates target account does not exist or does not belong to user.
	ErrAccountNotFound = errors.New("akun tidak ditemukan atau bukan milik pengguna")
	// ErrCategoryNotValid indicates category is inaccessible or type does not match.
	ErrCategoryNotValid = errors.New("kategori tidak valid untuk tipe transaksi ini")
)

// Transaction represents an income or expense record.
type Transaction struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id"`
	AccountID            string     `json:"account_id"`
	CategoryID           *string    `json:"category_id,omitempty"`
	Type                 string     `json:"type"` // 'income' or 'expense'
	Amount               int64      `json:"amount"`
	Currency             string     `json:"currency"`
	Description          string     `json:"description"`
	TransactionDate      time.Time  `json:"transaction_date"`
	Notes                *string    `json:"notes,omitempty"`
	Source               string     `json:"source"`
	CategorySource       string     `json:"category_source"`
	CategoryModelVersion *string    `json:"category_model_version,omitempty"`
	IdempotencyKey       *string    `json:"idempotency_key,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// TransactionDetailDTO represents an enriched transaction response with account & category metadata.
type TransactionDetailDTO struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	AccountName     string  `json:"account_name"`
	CategoryID      *string `json:"category_id,omitempty"`
	CategoryName    *string `json:"category_name,omitempty"`
	CategoryIcon    *string `json:"category_icon,omitempty"`
	CategoryColor   *string `json:"category_color,omitempty"`
	Type            string  `json:"type"`
	Amount          int64   `json:"amount"`
	Currency        string  `json:"currency"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	Notes           *string `json:"notes,omitempty"`
	Source          string  `json:"source"`
	CategorySource  string  `json:"category_source"`
	CreatedAt       string  `json:"created_at"`
}
