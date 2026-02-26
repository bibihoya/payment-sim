package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWallet(balance int64) *Wallet {
	now := time.Now()

	return &Wallet{
		ID:        uuid.New(),
		Balance:   balance,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Wallet) IsZero() bool {
	return w == nil || w.ID == uuid.Nil
}

func (w *Wallet) CanSend(amount int64) bool {
	if w.IsZero() {
		return false
	}
	return w.Balance >= amount
}
