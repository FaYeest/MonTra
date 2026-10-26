package transactions

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Filter specifies query criteria for listing transactions.
type Filter struct {
	AccountID  string
	CategoryID string
	Type       string
	StartDate  *time.Time
	EndDate    *time.Time
	Search     string
	Limit      int
	Offset     int
}

// Repository defines data access operations for transactions.
type Repository interface {
	Create(ctx context.Context, t *Transaction) error
	GetByID(ctx context.Context, userID, transactionID string) (*TransactionDetailDTO, error)
	GetByIdempotencyKey(ctx context.Context, userID, idempotencyKey string) (*TransactionDetailDTO, error)
	List(ctx context.Context, userID string, f Filter) ([]TransactionDetailDTO, int, error)
	Update(ctx context.Context, userID, transactionID string, t *Transaction) error
	Delete(ctx context.Context, userID, transactionID string) error
	ValidateAccount(ctx context.Context, userID, accountID string) error
	ValidateCategory(ctx context.Context, userID, categoryID, txType string) error
}

// PgRepository implements Repository using PostgreSQL.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates a new PostgreSQL transaction repository.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

// ValidateAccount checks whether the given account exists and belongs to the user.
func (r *PgRepository) ValidateAccount(ctx context.Context, userID, accountID string) error {
	query := `SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`
	var exists int
	err := r.pool.QueryRow(ctx, query, accountID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}
	return nil
}

// ValidateCategory checks whether the given category is accessible to the user and matches the transaction type.
func (r *PgRepository) ValidateCategory(ctx context.Context, userID, categoryID, txType string) error {
	query := `
		SELECT 1 FROM categories 
		WHERE id = $1 
		  AND (user_id = $2 OR (is_system = TRUE AND user_id IS NULL))
		  AND is_archived = FALSE
		  AND transaction_type = $3
	`
	var exists int
	err := r.pool.QueryRow(ctx, query, categoryID, userID, txType).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotValid
		}
		return err
	}
	return nil
}

// Create inserts a new transaction.
func (r *PgRepository) Create(ctx context.Context, t *Transaction) error {
	query := `
		INSERT INTO transactions (
			user_id, account_id, category_id, type, amount, currency,
			description, transaction_date, notes, source, category_source,
			category_model_version, idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		t.UserID,
		t.AccountID,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Currency,
		t.Description,
		t.TransactionDate,
		t.Notes,
		t.Source,
		t.CategorySource,
		t.CategoryModelVersion,
		t.IdempotencyKey,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

// GetByID retrieves a transaction by ID joined with account & category details.
func (r *PgRepository) GetByID(ctx context.Context, userID, transactionID string) (*TransactionDetailDTO, error) {
	query := `
		SELECT 
			t.id, t.account_id, a.name AS account_name,
			t.category_id, c.name AS category_name, c.icon AS category_icon, c.color AS category_color,
			t.type, t.amount, t.currency, t.description, t.transaction_date, t.notes,
			t.source, t.category_source, t.created_at
		FROM transactions t
		INNER JOIN accounts a ON a.id = t.account_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.id = $1 AND t.user_id = $2
	`
	return r.scanDetailRow(r.pool.QueryRow(ctx, query, transactionID, userID))
}

// GetByIdempotencyKey retrieves an existing transaction by user_id and idempotency_key.
func (r *PgRepository) GetByIdempotencyKey(ctx context.Context, userID, idempotencyKey string) (*TransactionDetailDTO, error) {
	if idempotencyKey == "" {
		return nil, ErrTransactionNotFound
	}
	query := `
		SELECT 
			t.id, t.account_id, a.name AS account_name,
			t.category_id, c.name AS category_name, c.icon AS category_icon, c.color AS category_color,
			t.type, t.amount, t.currency, t.description, t.transaction_date, t.notes,
			t.source, t.category_source, t.created_at
		FROM transactions t
		INNER JOIN accounts a ON a.id = t.account_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1 AND t.idempotency_key = $2
	`
	return r.scanDetailRow(r.pool.QueryRow(ctx, query, userID, idempotencyKey))
}

// List retrieves paginated transactions according to filter parameters.
func (r *PgRepository) List(ctx context.Context, userID string, f Filter) ([]TransactionDetailDTO, int, error) {
	var conditions []string
	var args []any
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("t.user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if f.AccountID != "" {
		conditions = append(conditions, fmt.Sprintf("t.account_id = $%d", argIndex))
		args = append(args, f.AccountID)
		argIndex++
	}
	if f.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("t.category_id = $%d", argIndex))
		args = append(args, f.CategoryID)
		argIndex++
	}
	if f.Type != "" {
		conditions = append(conditions, fmt.Sprintf("t.type = $%d", argIndex))
		args = append(args, f.Type)
		argIndex++
	}
	if f.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date >= $%d", argIndex))
		args = append(args, *f.StartDate)
		argIndex++
	}
	if f.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.transaction_date <= $%d", argIndex))
		args = append(args, *f.EndDate)
		argIndex++
	}
	if strings.TrimSpace(f.Search) != "" {
		conditions = append(conditions, fmt.Sprintf("t.description ILIKE $%d", argIndex))
		args = append(args, "%"+strings.TrimSpace(f.Search)+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions t WHERE %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	selectQuery := fmt.Sprintf(`
		SELECT 
			t.id, t.account_id, a.name AS account_name,
			t.category_id, c.name AS category_name, c.icon AS category_icon, c.color AS category_color,
			t.type, t.amount, t.currency, t.description, t.transaction_date, t.notes,
			t.source, t.category_source, t.created_at
		FROM transactions t
		INNER JOIN accounts a ON a.id = t.account_id
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE %s
		ORDER BY t.transaction_date DESC, t.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []TransactionDetailDTO
	for rows.Next() {
		var dto TransactionDetailDTO
		var txDate, createdAt time.Time
		err := rows.Scan(
			&dto.ID,
			&dto.AccountID,
			&dto.AccountName,
			&dto.CategoryID,
			&dto.CategoryName,
			&dto.CategoryIcon,
			&dto.CategoryColor,
			&dto.Type,
			&dto.Amount,
			&dto.Currency,
			&dto.Description,
			&txDate,
			&dto.Notes,
			&dto.Source,
			&dto.CategorySource,
			&createdAt,
		)
		if err != nil {
			return nil, 0, err
		}
		dto.TransactionDate = txDate.UTC().Format(time.RFC3339)
		dto.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		results = append(results, dto)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update modifies an existing transaction owned by user.
func (r *PgRepository) Update(ctx context.Context, userID, transactionID string, t *Transaction) error {
	query := `
		UPDATE transactions
		SET account_id = $1,
		    category_id = $2,
		    type = $3,
		    amount = $4,
		    description = $5,
		    transaction_date = $6,
		    notes = $7,
		    category_source = $8,
		    updated_at = NOW()
		WHERE id = $9 AND user_id = $10
	`
	tag, err := r.pool.Exec(ctx, query,
		t.AccountID,
		t.CategoryID,
		t.Type,
		t.Amount,
		t.Description,
		t.TransactionDate,
		t.Notes,
		t.CategorySource,
		transactionID,
		userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}
	return nil
}

// Delete removes a transaction owned by user.
func (r *PgRepository) Delete(ctx context.Context, userID, transactionID string) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	tag, err := r.pool.Exec(ctx, query, transactionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}
	return nil
}

func (r *PgRepository) scanDetailRow(row pgx.Row) (*TransactionDetailDTO, error) {
	var dto TransactionDetailDTO
	var txDate, createdAt time.Time
	err := row.Scan(
		&dto.ID,
		&dto.AccountID,
		&dto.AccountName,
		&dto.CategoryID,
		&dto.CategoryName,
		&dto.CategoryIcon,
		&dto.CategoryColor,
		&dto.Type,
		&dto.Amount,
		&dto.Currency,
		&dto.Description,
		&txDate,
		&dto.Notes,
		&dto.Source,
		&dto.CategorySource,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	dto.TransactionDate = txDate.UTC().Format(time.RFC3339)
	dto.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return &dto, nil
}
