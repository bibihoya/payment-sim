package handlers

import (
	"encoding/json"
	"net/http"
	"payment-sim/internal/domain"
	"payment-sim/internal/dto"
	"payment-sim/internal/storage"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	storage *storage.TransStorage
}

func NewTransactionHandler(storage *storage.TransStorage) *TransactionHandler {
	return &TransactionHandler{storage: storage}
}

// CreateTransaction POST /api/transactions
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tr := &domain.Transaction{
		ID:          uuid.New(),
		FromWalID:   uuid.MustParse(req.FromWalID),
		ToWalID:     uuid.MustParse(req.ToWalID),
		Amount:      req.Amount,
		Description: req.Description,
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
	}

	if err := h.storage.StoreTransaction(r.Context(), tr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.CreateTransactionResponse{
		ID:        tr.ID.String(),
		Status:    tr.Status.String(),
		Amount:    tr.Amount,
		CreatedAt: tr.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetTransaction GET /api/transactions/{id}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tr, err := h.storage.LoadTransaction(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tr == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := dto.CreateTransactionResponse{
		ID:        tr.ID.String(),
		Status:    tr.Status.String(),
		Amount:    tr.Amount,
		CreatedAt: tr.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
