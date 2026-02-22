package domain

import (
	"time"
)

type TransStatus int

const (
	StatusPending  TransStatus = iota // 0
	StatusApproved                    // 1
	StatusRejected                    // 2
	StatusFraud                       // 3
)

type Transaction struct {
	ID          int64
	FromWalID   int64
	ToWalID     int64
	Amount      int64
	Description string
	Status      TransStatus
	CreatedAt   time.Time
}
