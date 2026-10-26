package transfers

import (
	"errors"
	"time"
)

var (
	// ErrTransferNotFound indicates transfer was not found or is not owned by user.
	ErrTransferNotFound = errors.New("transfer tidak ditemukan")
	// ErrSameAccountTransfer indicates source and destination accounts are identical.
	ErrSameAccountTransfer = errors.New("akun sumber dan tujuan transfer tidak boleh sama")
	// ErrInvalidAmount indicates transfer amount <= 0.
	ErrInvalidAmount = errors.New("nominal transfer harus lebih besar dari 0")
	// ErrInvalidFee indicates negative fee amount.
	ErrInvalidFee = errors.New("biaya transfer tidak boleh negatif")
	// ErrAccountNotFound indicates source or destination account is invalid or not owned by user.
	ErrAccountNotFound = errors.New("akun sumber atau tujuan tidak ditemukan atau bukan milik pengguna")
	// ErrInvalidDate indicates transfer_date format cannot be parsed.
	ErrInvalidDate = errors.New("format tanggal transfer tidak valid")
)

// Transfer represents an internal fund movement between two accounts owned by the same user.
type Transfer struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id"`
	SourceAccountID      string     `json:"source_account_id"`
	DestinationAccountID string     `json:"destination_account_id"`
	Amount               int64      `json:"amount"`
	Currency             string     `json:"currency"`
	FeeAmount            int64      `json:"fee_amount"`
	TransactionDate      time.Time  `json:"transaction_date"`
	Description          *string    `json:"description,omitempty"`
	IdempotencyKey       *string    `json:"idempotency_key,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// TransferDetailDTO represents an enriched transfer record with source & destination account names.
type TransferDetailDTO struct {
	ID                     string  `json:"id"`
	SourceAccountID        string  `json:"source_account_id"`
	SourceAccountName      string  `json:"source_account_name"`
	DestinationAccountID   string  `json:"destination_account_id"`
	DestinationAccountName string  `json:"destination_account_name"`
	Amount                 int64   `json:"amount"`
	Currency               string  `json:"currency"`
	FeeAmount              int64   `json:"fee_amount"`
	TransactionDate        string  `json:"transaction_date"`
	Description            *string `json:"description,omitempty"`
	CreatedAt              string  `json:"created_at"`
}
