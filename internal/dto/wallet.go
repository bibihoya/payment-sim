package dto

import "time"

type CreateWalletRequest struct {
	Balance float64 `json:"balance,omitempty"`
}

type CreateWalletResponse struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type BalanceResponse struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionInfo struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Counterparty string `json:"counterparty"`
	Amount       int64  `json:"amount"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
}

type WalletHistory struct {
	ID           string            `json:"id"`
	Transactions []TransactionInfo `json:"transactions"`
	Total        int64             `json:"total"`
	Limit        int64             `json:"limit"`
}
