package worker

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"payment-sim/internal/kafka"
	"payment-sim/internal/storage"
	"payment-sim/internal/worker"
	"syscall"
	"time"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	transStorage := storage.NewTransStorage(db)
	walStorage := storage.NewWalStorage(db)

	brokers := []string{"localhost:9092"}
	consumer := kafka.NewConsumer(brokers, "transactions", "transaction-consumer")

	wk := worker.NewWorker(consumer, transStorage, walStorage)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := wk.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("error starting worker: %v", err)
		}
	}()

	log.Println("worker started")
	<-sigChan
	log.Println("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := consumer.Close(); err != nil {
		log.Printf("error closing consumer: %v", err)
	}

	<-shutdownCtx.Done()
	log.Println("worker stopped")
}
