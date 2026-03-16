package kafka

import (
	"encoding/json"
	"time"
)

type TransactionEvent struct {
	TransactionID string    `json:"transaction_id"`
	FromWalletID  string    `json:"from_wallet_id"`
	ToWalletID    string    `json:"to_wallet_id"`
	Amount        int64     `json:"amount"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

func (e *TransactionEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func TransEventFromJSON(data []byte) (*TransactionEvent, error) {
	var e TransactionEvent
	err := json.Unmarshal(data, &e)
	return &e, err
}
