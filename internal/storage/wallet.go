package storage

import (
	"context"
	"database/sql"
	"errors"
	"payment-sim/internal/domain"
	"time"

	"github.com/google/uuid"
)

type WalStorage struct {
	db *sql.DB
}

func NewWalStorage(db *sql.DB) *WalStorage {
	return &WalStorage{db: db}
}

func (wst *WalStorage) Create(ctx context.Context, balance int64) (*domain.Wallet, error) {
	w := domain.NewWallet(balance)

	query := `
		INSERT INTO wallets (id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := wst.db.ExecContext(ctx, query, w.ID, w.Balance, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (wst *WalStorage) LoadWallet(ctx context.Context, id string) (*domain.Wallet, error) {
	var w domain.Wallet
	var idStr string

	query := `
		SELECT id, balance, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`
	row := wst.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&idStr, &w.Balance, &w.CreatedAt, &w.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	w.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (wst *WalStorage) UpdateBalance(ctx context.Context, id string, balance int64) error {
	query := `
		UPDATE wallets
		SET balance = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := wst.db.ExecContext(ctx, query, balance, time.Now(), id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("nothing changed")
	}
	return nil
}

func (wst *WalStorage) GetLastTransactions(ctx context.Context, id string, limit int) ([]*domain.Transaction, error) {
	query := `
        SELECT id, from_wal_id, to_wal_id, amount, status, description, created_at
        FROM transactions
        WHERE from_wal_id = $1 OR to_wal_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `

	rows, err := wst.db.QueryContext(ctx, query, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var tr domain.Transaction
		var idStr, fromStr, toStr, statusStr string // ← status как string

		err := rows.Scan(&idStr, &fromStr, &toStr, &tr.Amount, &statusStr, &tr.Description, &tr.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Парсим UUID
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

		// Конвертируем строку статуса в TransStatus
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

		transactions = append(transactions, &tr)
	}

	return transactions, nil
}
