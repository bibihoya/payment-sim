package worker

import (
	"context"
	"log"
	"payment-sim/internal/antifraud"
	"payment-sim/internal/domain"
	"payment-sim/internal/kafka"
	"payment-sim/internal/storage"
)

type Worker struct {
	consumer      *kafka.Consumer
	transStorage  *storage.TransStorage
	walStorage    *storage.WalStorage
	fraudDetector *antifraud.Detector
}

func NewWorker(
	consumer *kafka.Consumer,
	transStorage *storage.TransStorage,
	walStorage *storage.WalStorage, fraudDetector *antifraud.Detector) *Worker {
	return &Worker{
		consumer:      consumer,
		transStorage:  transStorage,
		walStorage:    walStorage,
		fraudDetector: fraudDetector,
	}
}

func (wk *Worker) updateStatus(ctx context.Context, id string, status domain.TransStatus, msg string) error {
	tr, err := wk.transStorage.LoadTransaction(ctx, id)
	if err != nil {
		return err
	}
	if tr == nil {
		return nil
	}

	tr.Status = status
	return wk.transStorage.UpdateStatus(ctx, id, status, msg)
}

func (wk *Worker) processTransaction(ctx context.Context, event *kafka.TransactionEvent) error {
	log.Printf("Processing transaction: %s", event.TransactionID)

	tr, err := wk.transStorage.LoadTransaction(ctx, event.TransactionID)
	if err != nil {
		return err
	}
	if tr == nil {
		log.Printf("Transaction %s does not exist", event.TransactionID)
		return nil
	}
	if tr.Status != domain.StatusPending {
		log.Println("Transaction is not pending")
		return nil
	}

	fraud, reason, err := wk.fraudDetector.Check(ctx, tr)
	if err != nil {
		log.Printf("fraud error: %s", err.Error())
		return err
	}
	if fraud {
		log.Printf("Transaction %s is marked as fraud: %s", event.TransactionID, reason)
		return wk.updateStatus(ctx, event.TransactionID, domain.StatusFraud, reason)
	}

	err = wk.walStorage.Transfer(ctx, tr.FromWalID.String(), tr.ToWalID.String(), tr.Amount)
	if err != nil {
		log.Printf("Error transferring transaction: %s", err)

		switch err.Error() {
		case "sender wallet not found":
			return wk.updateStatus(ctx, tr.ID.String(), domain.StatusFailed, "sender wallet not found")
		case "receiver wallet not found":
			return wk.updateStatus(ctx, tr.ID.String(), domain.StatusFailed, "receiver wallet not found")
		case "sender wallet amount is low":
			return wk.updateStatus(ctx, tr.ID.String(), domain.StatusRejected, "sender wallet balance is low")
		default:
			return wk.updateStatus(ctx, tr.ID.String(), domain.StatusFailed, err.Error())
		}
	}

	log.Printf("Transaction %s approved", event.TransactionID)
	return wk.updateStatus(ctx, tr.ID.String(), domain.StatusApproved, "")
}

func (wk *Worker) Start(ctx context.Context) error {
	log.Println("Starting worker")

	return wk.consumer.Consume(ctx, wk.processTransaction)
}
