package antifraud

import (
	"context"
	"payment-sim/internal/domain"
	"time"
)

type Rule interface {
	Name() string
	Check(ctx context.Context, tr *domain.Transaction) (bool, string, error)
}

type AmountRule struct {
	maxAmount int64
}

func NewAmountRule(maxAmount int64) *AmountRule {
	return &AmountRule{
		maxAmount: maxAmount,
	}
}

func (ar *AmountRule) Name() string {
	return "amount_limit"
}

func (ar *AmountRule) Check(ctx context.Context, tr *domain.Transaction) (bool, string, error) {
	if tr.Amount > ar.maxAmount {
		return true, "amount exceeds limit", nil
	}
	return false, "", nil
}

type WalletStorage interface {
	CountTransactions(ctx context.Context, walletID string, since time.Time) (int, error)
}

type FrequencyRule struct {
	walStorage WalletStorage
	timeWindow time.Duration
	maxCount   int
}

func NewFrequencyRule(storage WalletStorage, timeWindow time.Duration, maxCount int) *FrequencyRule {
	return &FrequencyRule{
		walStorage: storage,
		timeWindow: timeWindow,
		maxCount:   maxCount,
	}
}

func (fr *FrequencyRule) Name() string {
	return "frequency_limit"
}

func (fr *FrequencyRule) Check(ctx context.Context, tr *domain.Transaction) (bool, string, error) {
	since := time.Now().Add(-fr.timeWindow)

	cnt, err := fr.walStorage.CountTransactions(ctx, tr.FromWalID.String(), since)
	if err != nil {
		return false, "", err
	}

	if cnt > fr.maxCount {
		return true, "too many transactions", nil
	}
	return false, "", nil
}

type SelfTransRule struct{}

func NewSelfTransRule() *SelfTransRule {
	return &SelfTransRule{}
}

func (sr *SelfTransRule) Name() string {
	return "self_transfer"
}
func (sr *SelfTransRule) Check(ctx context.Context, tr *domain.Transaction) (bool, string, error) {
	if tr.FromWalID == tr.ToWalID {
		return true, "self_transfer", nil
	}
	return false, "", nil
}
