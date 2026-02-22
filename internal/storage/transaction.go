package storage

import (
	"context"
	"database/sql"

	"payment-sim/internal/domain"
)

type TransStorage struct {
	db *sql.DB
}

func NewTransStorage(db *sql.DB) *TransStorage {
	return &TransStorage{db: db}
}

func (st *TransStorage) StoreTransaction(ctx context.Context, tr *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, from_wal_id, to_wal_id, amount, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := st.db.ExecContext(ctx, query,
		tr.ID, tr.FromWalID, tr.ToWalID, tr.Amount, tr.Status, tr.Description)

	return err
}

func (st *TransStorage) LoadTransaction(ctx context.Context, id int64) (*domain.Transaction, error) {
	var tr domain.Transaction

	query := `
		SELECT id, from_wal_id, to_wal_id, amount, status, description, created_at
		FROM transactions
		WHERE id = $1
	`
	row := st.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&tr.ID, &tr.FromWalID, &tr.ToWalID, &tr.Amount, &tr.Status, &tr.Description, &tr.CreatedAt)

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
