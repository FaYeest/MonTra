package users

import (
	"errors"
	"time"
)

var (
	// ErrNotFound indicates the requested user record does not exist.
	ErrNotFound = errors.New("pengguna tidak ditemukan")
	// ErrEmailAlreadyExists indicates a unique constraint violation on email.
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
)

// User represents the authoritative user entity.
type User struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"`
	DisplayName       string    `json:"display_name"`
	PreferredCurrency string    `json:"preferred_currency"`
	Timezone          string    `json:"timezone"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ProfileDTO represents sanitized user information exposed to API clients.
type ProfileDTO struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	DisplayName       string `json:"display_name"`
	PreferredCurrency string `json:"preferred_currency"`
	Timezone          string `json:"timezone"`
	CreatedAt         string `json:"created_at"`
}

// ToDTO converts an internal User model into a safe ProfileDTO.
func (u *User) ToDTO() ProfileDTO {
	return ProfileDTO{
		ID:                u.ID,
		Email:             u.Email,
		DisplayName:       u.DisplayName,
		PreferredCurrency: u.PreferredCurrency,
		Timezone:          u.Timezone,
		CreatedAt:         u.CreatedAt.UTC().Format(time.RFC3339),
	}
}
