package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"payment-sim/internal/domain"

	"github.com/google/uuid"
)

type TransStorage struct {
	db *sql.DB
}

func NewTransStorage(db *sql.DB) *TransStorage {
	return &TransStorage{db: db}
}

func (st *TransStorage) StoreTransaction(ctx context.Context, tr *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, from_wal_id, to_wal_id, amount, status, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	statusStr := tr.Status.String()

	_, err := st.db.ExecContext(ctx, query,
		tr.ID, tr.FromWalID, tr.ToWalID, tr.Amount, statusStr, tr.Description, time.Now())

	return err
}

func (st *TransStorage) LoadTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	var tr domain.Transaction
	var idStr, fromStr, toStr string

	query := `
		SELECT id, from_wal_id, to_wal_id, amount, status, description, created_at
		FROM transactions
		WHERE id = $1
	`
	row := st.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&idStr, &fromStr, &toStr, &tr.Amount, &tr.Status, &tr.Description, &tr.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, nil
		}
	}

	tr.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	tr.FromWalID, err = uuid.Parse(fromStr)
	if err != nil {
		return nil, err
	}
	tr.ToWalID, err = uuid.Parse(toStr)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}

func (st *TransStorage) UpdateStatus(ctx context.Context, id int64, status domain.TransStatus) error {
	query := `
		UPDATE transactions SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := st.db.ExecContext(ctx, query, status, id)
	return err
}
