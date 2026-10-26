package transfers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines data access operations for transfers.
type Repository interface {
	CreateAtomic(ctx context.Context, t *Transfer) error
	GetByID(ctx context.Context, userID, transferID string) (*TransferDetailDTO, error)
	GetByIdempotencyKey(ctx context.Context, userID, key string) (*TransferDetailDTO, error)
	List(ctx context.Context, userID string, limit, offset int) ([]TransferDetailDTO, int, error)
	Delete(ctx context.Context, userID, transferID string) error
}

// PgRepository implements Repository using PostgreSQL transactions.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates a new PostgreSQL transfer repository.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

// CreateAtomic records the transfer within an atomic database transaction.
func (r *PgRepository) CreateAtomic(ctx context.Context, t *Transfer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// 1. Validate and row-lock both source and destination accounts to ensure they belong to user
	validateQuery := `
		SELECT id 
		FROM accounts 
		WHERE id IN ($1, $2) AND user_id = $3 AND archived_at IS NULL
		FOR UPDATE
	`
	rows, err := tx.Query(ctx, validateQuery, t.SourceAccountID, t.DestinationAccountID, t.UserID)
	if err != nil {
		return err
	}
	var count int
	for rows.Next() {
		count++
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if count != 2 {
		return ErrAccountNotFound
	}

	// 2. Insert transfer record
	insertQuery := `
		INSERT INTO transfers (
			user_id, source_account_id, destination_account_id,
			amount, currency, fee_amount, transaction_date,
			description, idempotency_key
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, insertQuery,
		t.UserID,
		t.SourceAccountID,
		t.DestinationAccountID,
		t.Amount,
		t.Currency,
		t.FeeAmount,
		t.TransactionDate,
		t.Description,
		t.IdempotencyKey,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetByID retrieves a transfer record joined with account names.
func (r *PgRepository) GetByID(ctx context.Context, userID, transferID string) (*TransferDetailDTO, error) {
	query := `
		SELECT 
			tr.id,
			tr.source_account_id, sa.name AS source_account_name,
			tr.destination_account_id, da.name AS destination_account_name,
			tr.amount, tr.currency, tr.fee_amount,
			tr.transaction_date, tr.description, tr.created_at
		FROM transfers tr
		INNER JOIN accounts sa ON sa.id = tr.source_account_id
		INNER JOIN accounts da ON da.id = tr.destination_account_id
		WHERE tr.id = $1 AND tr.user_id = $2
	`
	return r.scanDetailRow(r.pool.QueryRow(ctx, query, transferID, userID))
}

// GetByIdempotencyKey retrieves an existing transfer by user and idempotency_key.
func (r *PgRepository) GetByIdempotencyKey(ctx context.Context, userID, key string) (*TransferDetailDTO, error) {
	if key == "" {
		return nil, ErrTransferNotFound
	}
	query := `
		SELECT 
			tr.id,
			tr.source_account_id, sa.name AS source_account_name,
			tr.destination_account_id, da.name AS destination_account_name,
			tr.amount, tr.currency, tr.fee_amount,
			tr.transaction_date, tr.description, tr.created_at
		FROM transfers tr
		INNER JOIN accounts sa ON sa.id = tr.source_account_id
		INNER JOIN accounts da ON da.id = tr.destination_account_id
		WHERE tr.user_id = $1 AND tr.idempotency_key = $2
	`
	return r.scanDetailRow(r.pool.QueryRow(ctx, query, userID, key))
}

// List returns paginated transfers for the user.
func (r *PgRepository) List(ctx context.Context, userID string, limit, offset int) ([]TransferDetailDTO, int, error) {
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	countQuery := `SELECT COUNT(*) FROM transfers WHERE user_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	selectQuery := `
		SELECT 
			tr.id,
			tr.source_account_id, sa.name AS source_account_name,
			tr.destination_account_id, da.name AS destination_account_name,
			tr.amount, tr.currency, tr.fee_amount,
			tr.transaction_date, tr.description, tr.created_at
		FROM transfers tr
		INNER JOIN accounts sa ON sa.id = tr.source_account_id
		INNER JOIN accounts da ON da.id = tr.destination_account_id
		WHERE tr.user_id = $1
		ORDER BY tr.transaction_date DESC, tr.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, selectQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []TransferDetailDTO
	for rows.Next() {
		var dto TransferDetailDTO
		var txDate, createdAt time.Time
		err := rows.Scan(
			&dto.ID,
			&dto.SourceAccountID,
			&dto.SourceAccountName,
			&dto.DestinationAccountID,
			&dto.DestinationAccountName,
			&dto.Amount,
			&dto.Currency,
			&dto.FeeAmount,
			&txDate,
			&dto.Description,
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

// Delete removes a transfer record within a transaction.
func (r *PgRepository) Delete(ctx context.Context, userID, transferID string) error {
	query := `DELETE FROM transfers WHERE id = $1 AND user_id = $2`
	tag, err := r.pool.Exec(ctx, query, transferID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTransferNotFound
	}
	return nil
}

func (r *PgRepository) scanDetailRow(row pgx.Row) (*TransferDetailDTO, error) {
	var dto TransferDetailDTO
	var txDate, createdAt time.Time
	err := row.Scan(
		&dto.ID,
		&dto.SourceAccountID,
		&dto.SourceAccountName,
		&dto.DestinationAccountID,
		&dto.DestinationAccountName,
		&dto.Amount,
		&dto.Currency,
		&dto.FeeAmount,
		&txDate,
		&dto.Description,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
		return nil, fmt.Errorf("failed to scan transfer row: %w", err)
	}
	dto.TransactionDate = txDate.UTC().Format(time.RFC3339)
	dto.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return &dto, nil
}
