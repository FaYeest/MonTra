package categories

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines database operations for categories.
type Repository interface {
	List(ctx context.Context, userID string, txType string) ([]Category, error)
	GetByID(ctx context.Context, userID, categoryID string) (*Category, error)
	Create(ctx context.Context, c *Category) error
}

// PgRepository implements Repository using PostgreSQL.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates a new PostgreSQL category repository.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

// List returns all active categories accessible to the user (system categories + user-created categories),
// optionally filtered by transaction type ('income' or 'expense').
func (r *PgRepository) List(ctx context.Context, userID string, txType string) ([]Category, error) {
	var query string
	var rows pgx.Rows
	var err error

	txType = strings.ToLower(strings.TrimSpace(txType))
	if txType != "" {
		query = `
			SELECT id, user_id, name, transaction_type, COALESCE(icon, ''), COALESCE(color, ''), is_system, is_archived, created_at, updated_at
			FROM categories
			WHERE (user_id = $1 OR (is_system = TRUE AND user_id IS NULL))
			  AND is_archived = FALSE
			  AND transaction_type = $2
			ORDER BY is_system DESC, name ASC
		`
		rows, err = r.pool.Query(ctx, query, userID, txType)
	} else {
		query = `
			SELECT id, user_id, name, transaction_type, COALESCE(icon, ''), COALESCE(color, ''), is_system, is_archived, created_at, updated_at
			FROM categories
			WHERE (user_id = $1 OR (is_system = TRUE AND user_id IS NULL))
			  AND is_archived = FALSE
			ORDER BY is_system DESC, name ASC
		`
		rows, err = r.pool.Query(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.TransactionType,
			&c.Icon,
			&c.Color,
			&c.IsSystem,
			&c.IsArchived,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetByID retrieves a single category by ID accessible to the user.
func (r *PgRepository) GetByID(ctx context.Context, userID, categoryID string) (*Category, error) {
	query := `
		SELECT id, user_id, name, transaction_type, COALESCE(icon, ''), COALESCE(color, ''), is_system, is_archived, created_at, updated_at
		FROM categories
		WHERE id = $1
		  AND (user_id = $2 OR (is_system = TRUE AND user_id IS NULL))
		  AND is_archived = FALSE
	`
	var c Category
	err := r.pool.QueryRow(ctx, query, categoryID, userID).Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.TransactionType,
		&c.Icon,
		&c.Color,
		&c.IsSystem,
		&c.IsArchived,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

// Create inserts a custom category for the given user.
func (r *PgRepository) Create(ctx context.Context, c *Category) error {
	query := `
		INSERT INTO categories (user_id, name, transaction_type, icon, color, is_system, is_archived)
		VALUES ($1, $2, $3, $4, $5, FALSE, FALSE)
		RETURNING id, is_system, is_archived, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		c.UserID,
		strings.TrimSpace(c.Name),
		strings.ToLower(strings.TrimSpace(c.TransactionType)),
		strings.TrimSpace(c.Icon),
		strings.TrimSpace(c.Color),
	).Scan(
		&c.ID,
		&c.IsSystem,
		&c.IsArchived,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
}
