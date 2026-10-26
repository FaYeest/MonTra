package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken indicates the JWT signature or claims are invalid.
	ErrInvalidToken = errors.New("token tidak valid atau telah kedaluwarsa")
)

// Claims represents the JWT payload for MonTra users.
type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenManager manages JWT signing and verification.
type TokenManager struct {
	secret []byte
	expiry time.Duration
}

// NewTokenManager initializes a TokenManager with a secret key and default token expiry.
func NewTokenManager(secret string, expiry time.Duration) (*TokenManager, error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secret cannot be empty")
	}
	if expiry == 0 {
		expiry = 72 * time.Hour
	}
	return &TokenManager{
		secret: []byte(secret),
		expiry: expiry,
	}, nil
}

// Generate issues a signed JWT for the given user ID and email.
func (tm *TokenManager) Generate(userID, email string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.expiry)),
			Issuer:    "montra-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

// Verify validates the JWT string and extracts the claims.
func (tm *TokenManager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return tm.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.UserID == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
