package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	FromWalID   string `json:"from_wal_id"`
	ToWalID     string `json:"to_wal_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

func (r *CreateTransactionRequest) Validate() error {
	if r.Amount < 0 {
		return errors.New("amount must be positive")
	}

	if _, err := uuid.Parse(r.ToWalID); err != nil {
		return errors.New("to_wal_id must be uuid")
	}

	if _, err := uuid.Parse(r.FromWalID); err != nil {
		return errors.New("from_wal_id must be uuid")
	}

	if r.ToWalID == r.FromWalID {
		return errors.New("to_wal_id must not be the same as from_wal_id")
	}

	return nil
}

type CreateTransactionResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
