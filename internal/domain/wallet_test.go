package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestWallet_CanSend(t *testing.T) {
	tests := []struct {
		name     string
		balance  int64
		amount   int64
		expected bool
	}{
		{
			name:     "достаточно денег",
			balance:  1000,
			amount:   500,
			expected: true,
		},
		{
			name:     "ровно столько же",
			balance:  1000,
			amount:   1000,
			expected: true,
		},
		{
			name:     "недостаточно денег",
			balance:  1000,
			amount:   1500,
			expected: false,
		},
		{
			name:     "нулевой баланс",
			balance:  0,
			amount:   100,
			expected: false,
		},
		{
			name:     "отрицательная сумма (не должна проходить)",
			balance:  1000,
			amount:   -100,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := &Wallet{
				ID:        uuid.New(),
				Balance:   tt.balance,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			result := wallet.CanSend(tt.amount)
			if result != tt.expected {
				t.Errorf("CanSend() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWallet_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		wallet   *Wallet
		expected bool
	}{
		{
			name:     "nil кошелек",
			wallet:   nil,
			expected: true,
		},
		{
			name:     "пустой ID",
			wallet:   &Wallet{ID: uuid.Nil},
			expected: true,
		},
		{
			name:     "нормальный кошелек",
			wallet:   &Wallet{ID: uuid.New()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wallet.IsZero(); got != tt.expected {
				t.Errorf("IsZero() = %v, want %v", got, tt.expected)
			}
		})
	}
}
