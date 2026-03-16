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
		INSERT INTO transactions (id, from_wal_id, to_wal_id, amount, status, description, created_at, updated_at, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	statusStr := tr.Status.String()

	_, err := st.db.ExecContext(ctx, query,
		tr.ID, tr.FromWalID, tr.ToWalID, tr.Amount, statusStr, tr.Description, time.Now(), time.Now(), tr.ErrorMsg)

	return err
}

func (st *TransStorage) LoadTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	var tr domain.Transaction
	var idStr, fromStr, toStr, statusStr string
	var errMsg sql.NullString

	query := `
		SELECT id, from_wal_id, to_wal_id, amount, status, description, created_at, error_message
		FROM transactions
		WHERE id = $1
	`
	row := st.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&idStr, &fromStr, &toStr, &tr.Amount, &statusStr, &tr.Description, &tr.CreatedAt, &errMsg)

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

	switch statusStr {
	case "pending":
		tr.Status = domain.StatusPending
	case "approved":
		tr.Status = domain.StatusApproved
	case "rejected":
		tr.Status = domain.StatusRejected
	case "fraud":
		tr.Status = domain.StatusFraud
	default:
		tr.Status = domain.StatusPending
	}

	if errMsg.Valid {
		tr.ErrorMsg = errMsg.String
	}

	return &tr, nil
}

func (st *TransStorage) UpdateStatus(ctx context.Context, id string, status domain.TransStatus, errorMsg string) error {
	query := `
        UPDATE transactions 
        SET status = $1, 
            error_message = $2, 
            updated_at = NOW() 
        WHERE id = $3
    `
	_, err := st.db.ExecContext(ctx, query, status.String(), errorMsg, id)
	return err
}
