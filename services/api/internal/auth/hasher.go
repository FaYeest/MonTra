package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooShort indicates password does not meet minimum length requirement.
	ErrPasswordTooShort = errors.New("kata sandi minimal 8 karakter")
	// ErrPasswordTooLong indicates password exceeds bcrypt 72-byte limit.
	ErrPasswordTooLong = errors.New("kata sandi maksimal 72 byte")
)

// Hasher handles secure password hashing and comparison.
type Hasher struct {
	cost int
}

// NewHasher creates a Hasher with standard bcrypt cost.
func NewHasher() *Hasher {
	return &Hasher{cost: bcrypt.DefaultCost}
}

// Hash generates a bcrypt hash of the password.
func (h *Hasher) Hash(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}
	if len([]byte(password)) > 72 {
		return "", ErrPasswordTooLong
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare verifies whether the plain password matches the bcrypt hash.
func (h *Hasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
