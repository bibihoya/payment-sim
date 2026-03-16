package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-sim/internal/handlers"
	"payment-sim/internal/kafka"
	"payment-sim/internal/storage"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	brokers := []string{"localhost:9092"}
	producer := kafka.NewProducer(brokers, "transactions")
	defer producer.Close()

	transStorage := storage.NewTransStorage(db)
	walStorage := storage.NewWalStorage(db)

	transHandler := handlers.NewTransactionHandler(transStorage, producer)
	walletHandler := handlers.NewWalletHandler(walStorage)

	r := mux.NewRouter()
	r.HandleFunc("/api/transactions", transHandler.CreateTransaction).Methods("POST")
	r.HandleFunc("/api/transactions/{id}", transHandler.GetTransaction).Methods("GET")

	r.HandleFunc("/api/wallets", walletHandler.CreateWallet).Methods("POST")
	r.HandleFunc("/api/wallets/{id}/balance", walletHandler.GetBalance).Methods("GET")
	r.HandleFunc("/api/wallets/{id}/transactions", walletHandler.GetWalletHistory).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Канал для сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down API...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if err := producer.Close(); err != nil {
		log.Printf("Producer close error: %v", err)
	}

	log.Println("API stopped")
}
