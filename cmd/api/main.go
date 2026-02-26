package main

import (
	"database/sql"
	"log"
	"net/http"
	"payment-sim/internal/handlers"
	"payment-sim/internal/storage"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	transStorage := storage.NewTransStorage(db)
	walStorage := storage.NewWalStorage(db)

	transHandler := handlers.NewTransactionHandler(transStorage)
	walletHandler := handlers.NewWalletHandler(walStorage)

	r := mux.NewRouter()
	r.HandleFunc("/api/transactions", transHandler.CreateTransaction).Methods("POST")
	r.HandleFunc("/api/transactions/{id}", transHandler.GetTransaction).Methods("GET")

	r.HandleFunc("/api/wallets", walletHandler.CreateWallet).Methods("POST")
	r.HandleFunc("/api/wallets/{id}/balance", walletHandler.GetBalance).Methods("GET")
	r.HandleFunc("/api/wallets/{id}/transactions", walletHandler.GetWalletHistory).Methods("GET")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
