package auth

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/montra/montra/services/api/internal/users"
)

var (
	// ErrInvalidEmail indicates malformed or oversized email.
	ErrInvalidEmail = errors.New("format email tidak valid")
	// ErrInvalidDisplayName indicates invalid display name length.
	ErrInvalidDisplayName = errors.New("nama tampilan harus antara 1 sampai 60 karakter")
	// ErrInvalidCredentials indicates incorrect email or password.
	ErrInvalidCredentials = errors.New("email atau kata sandi tidak valid")
)

// RegisterRequest holds registration inputs.
type RegisterRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	DisplayName       string `json:"display_name"`
	FullName          string `json:"full_name,omitempty"`
	PreferredCurrency string `json:"preferred_currency,omitempty"`
}

// LoginRequest holds authentication credentials.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse returns the JWT bearer token and public user profile.
type AuthResponse struct {
	Token       string           `json:"token"`
	AccessToken string           `json:"access_token,omitempty"`
	User        users.ProfileDTO `json:"user"`
}

// Service provides authentication and profile query logic.
type Service struct {
	userRepo users.Repository
	hasher   *Hasher
	tm       *TokenManager
}

// NewService constructs an auth Service.
func NewService(userRepo users.Repository, hasher *Hasher, tm *TokenManager) *Service {
	return &Service{
		userRepo: userRepo,
		hasher:   hasher,
		tm:       tm,
	}
}

// Register creates a new user, hashes password, saves record, and issues a JWT token.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || len(email) > 128 {
		return nil, ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" && req.FullName != "" {
		displayName = strings.TrimSpace(req.FullName)
	}
	if len(displayName) < 1 || len(displayName) > 60 {
		return nil, ErrInvalidDisplayName
	}

	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(req.PreferredCurrency))
	if currency == "" {
		currency = "IDR"
	}
	if len(currency) != 3 {
		return nil, errors.New("preferred_currency harus 3 karakter ISO 4217")
	}

	u := &users.User{
		Email:             email,
		PasswordHash:      hash,
		DisplayName:       displayName,
		PreferredCurrency: currency,
		Timezone:          "Asia/Jakarta",
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	token, err := s.tm.Generate(u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:       token,
		AccessToken: token,
		User:        u.ToDTO(),
	}, nil
}

// Login verifies credentials and issues a JWT token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !s.hasher.Compare(u.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tm.Generate(u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:       token,
		AccessToken: token,
		User:        u.ToDTO(),
	}, nil
}

// GetProfile retrieves the profile of the authenticated user.
func (s *Service) GetProfile(ctx context.Context, userID string) (*users.ProfileDTO, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	dto := u.ToDTO()
	return &dto, nil
}
