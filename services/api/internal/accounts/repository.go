package accounts

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines data access methods for accounts.
type Repository interface {
	ListWithBalance(ctx context.Context, userID string) ([]Account, error)
	GetByIDWithBalance(ctx context.Context, userID, accountID string) (*Account, error)
	Create(ctx context.Context, a *Account) error
	Update(ctx context.Context, userID, accountID, name, accountType string) error
	Archive(ctx context.Context, userID, accountID string) error
}

// PgRepository implements Repository using PostgreSQL via pgxpool.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates a new PostgreSQL account repository.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

// ListWithBalance returns all non-archived accounts for the user with deterministically calculated running balances.
func (r *PgRepository) ListWithBalance(ctx context.Context, userID string) ([]Account, error) {
	query := `
		SELECT 
			a.id, 
			a.user_id, 
			a.name, 
			a.account_type, 
			a.opening_balance, 
			a.currency, 
			a.archived_at, 
			a.created_at, 
			a.updated_at,
			(
				a.opening_balance 
				+ COALESCE((SELECT SUM(amount) FROM transactions WHERE account_id = a.id AND type = 'income'), 0)
				- COALESCE((SELECT SUM(amount) FROM transactions WHERE account_id = a.id AND type = 'expense'), 0)
				+ COALESCE((SELECT SUM(amount) FROM transfers WHERE destination_account_id = a.id), 0)
				- COALESCE((SELECT SUM(amount + fee_amount) FROM transfers WHERE source_account_id = a.id), 0)
			) AS current_balance
		FROM accounts a
		WHERE a.user_id = $1 AND a.archived_at IS NULL
		ORDER BY a.created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Account
	for rows.Next() {
		var a Account
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Name,
			&a.AccountType,
			&a.OpeningBalance,
			&a.Currency,
			&a.ArchivedAt,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.CurrentBalance,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByIDWithBalance returns an account by ID with its running balance, ensuring strict user_id scoping.
func (r *PgRepository) GetByIDWithBalance(ctx context.Context, userID, accountID string) (*Account, error) {
	query := `
		SELECT 
			a.id, 
			a.user_id, 
			a.name, 
			a.account_type, 
			a.opening_balance, 
			a.currency, 
			a.archived_at, 
			a.created_at, 
			a.updated_at,
			(
				a.opening_balance 
				+ COALESCE((SELECT SUM(amount) FROM transactions WHERE account_id = a.id AND type = 'income'), 0)
				- COALESCE((SELECT SUM(amount) FROM transactions WHERE account_id = a.id AND type = 'expense'), 0)
				+ COALESCE((SELECT SUM(amount) FROM transfers WHERE destination_account_id = a.id), 0)
				- COALESCE((SELECT SUM(amount + fee_amount) FROM transfers WHERE source_account_id = a.id), 0)
			) AS current_balance
		FROM accounts a
		WHERE a.id = $1 AND a.user_id = $2
	`
	var a Account
	err := r.pool.QueryRow(ctx, query, accountID, userID).Scan(
		&a.ID,
		&a.UserID,
		&a.Name,
		&a.AccountType,
		&a.OpeningBalance,
		&a.Currency,
		&a.ArchivedAt,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.CurrentBalance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &a, nil
}

// Create inserts a new account record.
func (r *PgRepository) Create(ctx context.Context, a *Account) error {
	query := `
		INSERT INTO accounts (user_id, name, account_type, opening_balance, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		a.UserID,
		strings.TrimSpace(a.Name),
		strings.ToLower(strings.TrimSpace(a.AccountType)),
		a.OpeningBalance,
		a.Currency,
	).Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return err
	}
	a.CurrentBalance = a.OpeningBalance
	return nil
}

// Update modifies the name or account_type of an existing account owned by the user.
func (r *PgRepository) Update(ctx context.Context, userID, accountID, name, accountType string) error {
	query := `
		UPDATE accounts
		SET name = COALESCE(NULLIF($1, ''), name),
		    account_type = COALESCE(NULLIF($2, ''), account_type),
		    updated_at = NOW()
		WHERE id = $3 AND user_id = $4 AND archived_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, query,
		strings.TrimSpace(name),
		strings.ToLower(strings.TrimSpace(accountType)),
		accountID,
		userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrAccountNotFound
	}
	return nil
}

// Archive sets archived_at timestamp on the account so it is hidden from active balance lists without losing history.
func (r *PgRepository) Archive(ctx context.Context, userID, accountID string) error {
	query := `
		UPDATE accounts
		SET archived_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND archived_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, query, accountID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrAccountNotFound
	}
	return nil
}
