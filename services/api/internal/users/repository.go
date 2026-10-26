package users

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines storage operations for users.
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

// PgRepository implements Repository using PostgreSQL via pgxpool.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates a new PostgreSQL-backed user repository.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

// Create inserts a new user record.
func (r *PgRepository) Create(ctx context.Context, u *User) error {
	normalizedEmail := strings.ToLower(strings.TrimSpace(u.Email))
	query := `
		INSERT INTO users (email, password_hash, display_name, preferred_currency, timezone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, display_name, preferred_currency, timezone, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		normalizedEmail,
		u.PasswordHash,
		strings.TrimSpace(u.DisplayName),
		u.PreferredCurrency,
		u.Timezone,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.DisplayName,
		&u.PreferredCurrency,
		&u.Timezone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

// GetByEmail retrieves a user by their case-insensitive normalized email address.
func (r *PgRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	query := `
		SELECT id, email, password_hash, display_name, preferred_currency, timezone, created_at, updated_at
		FROM users
		WHERE LOWER(email) = $1
	`
	var u User
	err := r.pool.QueryRow(ctx, query, normalizedEmail).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.DisplayName,
		&u.PreferredCurrency,
		&u.Timezone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetByID retrieves a user by their UUID primary key.
func (r *PgRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, display_name, preferred_currency, timezone, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.DisplayName,
		&u.PreferredCurrency,
		&u.Timezone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
