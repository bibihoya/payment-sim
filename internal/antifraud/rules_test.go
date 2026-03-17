package antifraud

import (
	"context"
	"payment-sim/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockWalletStorage struct {
	count int
	err   error
}

func (m *mockWalletStorage) CountTransactions(ctx context.Context, walletID string, since time.Time) (int, error) {
	return m.count, m.err
}

func TestAmountRule_Check(t *testing.T) {
	rule := NewAmountRule(1000000)

	tests := []struct {
		name   string
		amount int64
		fraud  bool
		reason string
	}{
		{
			name:   "меньше лимита",
			amount: 500000,
			fraud:  false,
			reason: "",
		},
		{
			name:   "ровно лимит",
			amount: 1000000,
			fraud:  false,
			reason: "",
		},
		{
			name:   "больше лимита",
			amount: 2000000,
			fraud:  true,
			reason: "amount exceeds limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &domain.Transaction{Amount: tt.amount}
			fraud, reason, err := rule.Check(context.Background(), tr)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if fraud != tt.fraud {
				t.Errorf("fraud = %v, want %v", fraud, tt.fraud)
			}
			if reason != tt.reason {
				t.Errorf("reason = %q, want %q", reason, tt.reason)
			}
		})
	}
}

func TestFrequencyRule_Check(t *testing.T) {
	mock := &mockWalletStorage{count: 5}
	rule := NewFrequencyRule(mock, 5*time.Minute, 10)

	tr := &domain.Transaction{
		FromWalID: uuid.New(),
	}

	mock.count = 5
	fraud, reason, err := rule.Check(context.Background(), tr)
	if err != nil || fraud || reason != "" {
		t.Errorf("expected no fraud, got fraud=%v, reason=%q, err=%v", fraud, reason, err)
	}

	mock.count = 10
	fraud, reason, err = rule.Check(context.Background(), tr)
	if err != nil || fraud || reason != "" {
		t.Errorf("expected no fraud at limit, got fraud=%v", fraud)
	}

	mock.count = 11
	fraud, reason, err = rule.Check(context.Background(), tr)
	if err != nil || !fraud || reason != "too many transactions" {
		t.Errorf("expected fraud, got fraud=%v, reason=%q", fraud, reason)
	}
}

func TestSelfTransferRule_Check(t *testing.T) {
	rule := &SelfTransRule{}
	walletID := uuid.New()

	tests := []struct {
		name   string
		from   uuid.UUID
		to     uuid.UUID
		fraud  bool
		reason string
	}{
		{
			name:  "разные кошельки",
			from:  uuid.New(),
			to:    uuid.New(),
			fraud: false,
		},
		{
			name:   "себе",
			from:   walletID,
			to:     walletID,
			fraud:  true,
			reason: "self_transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &domain.Transaction{
				FromWalID: tt.from,
				ToWalID:   tt.to,
			}
			fraud, reason, err := rule.Check(context.Background(), tr)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if fraud != tt.fraud {
				t.Errorf("fraud = %v, want %v", fraud, tt.fraud)
			}
			if reason != tt.reason {
				t.Errorf("reason = %q, want %q", reason, tt.reason)
			}
		})
	}
}
