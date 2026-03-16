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
	StatusFailed                      // 3
	StatusFraud                       // 4
)

type Transaction struct {
	ID          uuid.UUID
	FromWalID   uuid.UUID
	ToWalID     uuid.UUID
	Amount      int64
	Description string
	Status      TransStatus
	CreatedAt   time.Time
	ErrorMsg    string
}

func (s TransStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusApproved:
		return "approved"
	case StatusRejected:
		return "rejected"
	case StatusFailed:
		return "failed"
	case StatusFraud:
		return "fraud"
	default:
		return "unknown"
	}
}
