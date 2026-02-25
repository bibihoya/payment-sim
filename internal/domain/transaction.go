package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransStatus int

const (
	StatusPending  TransStatus = iota // 0
	StatusApproved                    // 1
	StatusRejected                    // 2
	StatusFraud                       // 3
)

type Transaction struct {
	ID          uuid.UUID
	FromWalID   uuid.UUID
	ToWalID     uuid.UUID
	Amount      int64
	Description string
	Status      TransStatus
	CreatedAt   time.Time
}
