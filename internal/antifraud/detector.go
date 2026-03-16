package antifraud

import (
	"context"
	"log"
	"payment-sim/internal/domain"
	"payment-sim/internal/storage"
	"time"
)

type Detector struct {
	rules []Rule
}

func NewDetector(wst *storage.WalStorage) *Detector {
	return &Detector{
		rules: []Rule{
			NewAmountRule(1000000),
			NewFrequencyRule(wst, 5*time.Minute, 10),
			NewSelfTransRule(),
		},
	}
}

func (d *Detector) Check(ctx context.Context, tr *domain.Transaction) (bool, string, error) {
	log.Printf("Running fraud detection for transaction %s", tr.ID.String())

	for _, rule := range d.rules {
		fraud, reason, err := rule.Check(ctx, tr)
		if err != nil {
			log.Printf("Rule %s error %s", rule, err.Error())
			continue
		}

		if fraud {
			log.Printf("Fraud detected, rule %s", rule)
			return true, reason, nil
		}
	}

	log.Printf("Transaction %s passed", tr.ID.String())
	return false, "", nil
}
