package integration

import (
	"context"
	"database/sql"
	"payment-sim/internal/domain"
	"testing"
	"time"

	"payment-sim/internal/storage"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTransferMoney(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer pgContainer.Terminate(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE wallets (
            id UUID PRIMARY KEY,
            balance BIGINT NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        );

        CREATE TABLE transactions (
            id UUID PRIMARY KEY,
            from_wal_id UUID NOT NULL REFERENCES wallets(id),
            to_wal_id UUID NOT NULL REFERENCES wallets(id),
            amount BIGINT NOT NULL,
            status TEXT NOT NULL,
            description TEXT,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL,
            error_message TEXT
        );
    `)
	if err != nil {
		t.Fatal(err)
	}

	walStorage := storage.NewWalStorage(db)

	fromWallet := domain.NewWallet(1000000)
	toWallet := domain.NewWallet(50000)

	_, err = db.Exec(
		`INSERT INTO wallets (id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		fromWallet.ID, fromWallet.Balance, fromWallet.CreatedAt, fromWallet.UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(
		`INSERT INTO wallets (id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		toWallet.ID, toWallet.Balance, toWallet.CreatedAt, toWallet.UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("successful transfer", func(t *testing.T) {
		err := walStorage.Transfer(ctx, fromWallet.ID.String(), toWallet.ID.String(), 30000)
		if err != nil {
			t.Fatalf("Transfer failed: %v", err)
		}

		// Проверяем балансы
		var fromBalance, toBalance int64
		err = db.QueryRow("SELECT balance FROM wallets WHERE id = $1", fromWallet.ID).Scan(&fromBalance)
		if err != nil || fromBalance != 1000000-30000 {
			t.Errorf("from wallet balance = %d, want %d", fromBalance, 1000000-30000)
		}

		err = db.QueryRow("SELECT balance FROM wallets WHERE id = $1", toWallet.ID).Scan(&toBalance)
		if err != nil || toBalance != 50000+30000 {
			t.Errorf("to wallet balance = %d, want %d", toBalance, 50000+30000)
		}
	})

	t.Run("insufficient funds", func(t *testing.T) {
		// Создаём бедный кошелёк
		poorWallet := domain.NewWallet(100)
		_, err = db.Exec(
			`INSERT INTO wallets (id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
			poorWallet.ID, poorWallet.Balance, poorWallet.CreatedAt, poorWallet.UpdatedAt,
		)
		if err != nil {
			t.Fatal(err)
		}

		err := walStorage.Transfer(ctx, poorWallet.ID.String(), toWallet.ID.String(), 200)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "sender wallet amount is low" {
			t.Errorf("expected 'sender wallet amount is low', got %v", err)
		}
	})

	t.Run("sender not found", func(t *testing.T) {
		err := walStorage.Transfer(ctx, "00000000-0000-0000-0000-000000000000", toWallet.ID.String(), 100)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "sender wallet not found" {
			t.Errorf("expected 'sender wallet not found', got %v", err)
		}
	})

	t.Run("receiver not found", func(t *testing.T) {
		err := walStorage.Transfer(ctx, fromWallet.ID.String(), "00000000-0000-0000-0000-000000000000", 100)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "receiver wallet not found" {
			t.Errorf("expected 'receiver wallet not found', got %v", err)
		}
	})

}
