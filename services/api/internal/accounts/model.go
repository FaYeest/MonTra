package accounts

import (
	"errors"
	"time"
)

var (
	// ErrAccountNotFound indicates the account does not exist or does not belong to the user.
	ErrAccountNotFound = errors.New("akun tidak ditemukan")
	// ErrInvalidAccountName indicates the account name length is invalid.
	ErrInvalidAccountName = errors.New("nama akun harus antara 1 sampai 50 karakter")
	// ErrInvalidAccountType indicates an unsupported account type.
	ErrInvalidAccountType = errors.New("tipe akun harus salah satu dari: 'cash', 'bank', 'ewallet', 'other'")
	// ErrNegativeOpeningBalance indicates opening balance cannot be negative.
	ErrNegativeOpeningBalance = errors.New("saldo awal tidak boleh negatif")
)

// ValidAccountTypes lists the supported account types.
var ValidAccountTypes = map[string]bool{
	"cash":    true,
	"bank":    true,
	"ewallet": true,
	"other":   true,
}

// Account represents a financial wallet or account with running balance.
type Account struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Name           string     `json:"name"`
	AccountType    string     `json:"account_type"`
	OpeningBalance int64      `json:"opening_balance"`
	CurrentBalance int64      `json:"current_balance"`
	Currency       string     `json:"currency"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// AccountDTO is the client-facing representation of an account.
type AccountDTO struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AccountType    string `json:"account_type"`
	OpeningBalance int64  `json:"opening_balance"`
	CurrentBalance int64  `json:"current_balance"`
	Currency       string `json:"currency"`
	CreatedAt      string `json:"created_at"`
}

// ToDTO converts an internal Account model into a client-safe AccountDTO.
func (a *Account) ToDTO() AccountDTO {
	return AccountDTO{
		ID:             a.ID,
		Name:           a.Name,
		AccountType:    a.AccountType,
		OpeningBalance: a.OpeningBalance,
		CurrentBalance: a.CurrentBalance,
		Currency:       a.Currency,
		CreatedAt:      a.CreatedAt.UTC().Format(time.RFC3339),
	}
}
