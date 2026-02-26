package handlers

import (
	"encoding/json"
	"net/http"
	"payment-sim/internal/dto"
	"payment-sim/internal/storage"
	"strconv"

	"github.com/gorilla/mux"
)

type WalletHandler struct {
	storage *storage.WalStorage
}

func NewWalletHandler(storage *storage.WalStorage) *WalletHandler {
	return &WalletHandler{storage: storage}
}

// CreateWallet POST /api/wallets
func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Balance < 0 {
		http.Error(w, "balance must be positive", http.StatusBadRequest)
		return
	}

	wal, err := h.storage.Create(r.Context(), req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.CreateWalletResponse{
		ID:        wal.ID.String(),
		Balance:   wal.Balance,
		CreatedAt: wal.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBalance GET /api/wallets/{id}/balance
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	wal, err := h.storage.LoadWallet(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if wal == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	resp := dto.BalanceResponse{
		ID:        wal.ID.String(),
		Balance:   wal.Balance,
		UpdatedAt: wal.UpdatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetWalletHistory — GET /api/wallets/{id}/transactions?limit=10
func (h *WalletHandler) GetWalletHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if lim, err := strconv.Atoi(limitStr); err == nil {
			limit = lim
		}
	}

	wal, err := h.storage.LoadWallet(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if wal == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	transactions, err := h.storage.GetLastTransactions(r.Context(), id, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []dto.TransactionInfo
	for _, tr := range transactions {
		item := dto.TransactionInfo{
			ID:          tr.ID.String(),
			Amount:      tr.Amount,
			Status:      tr.Status.String(),
			Description: tr.Description,
		}

		if tr.FromWalID.String() == id {
			item.Type = "outgoing"
			item.Counterparty = tr.ToWalID.String()
		} else {
			item.Type = "incoming"
			item.Counterparty = tr.FromWalID.String()
		}

		items = append(items, item)
	}

	resp := dto.WalletHistory{
		ID:           id,
		Transactions: items,
		Total:        int64(len(items)),
		Limit:        int64(limit),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
