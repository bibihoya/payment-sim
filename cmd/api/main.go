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
	transHandler := handlers.NewTransactionHandler(transStorage)

	r := mux.NewRouter()
	r.HandleFunc("/api/transactions", transHandler.CreateTransaction).Methods("POST")
	r.HandleFunc("/api/transactions/{id}", transHandler.GetTransaction).Methods("GET")

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
