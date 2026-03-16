package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Second,
		Async:        false,
		RequiredAcks: kafka.RequireAll,
		Compression:  kafka.Snappy,
	}

	log.Printf("Kafka producer created for topic: %s", topic)
	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) PublishTransaction(ctx context.Context, event *TransactionEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal transaction event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.TransactionID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte("transaction_created")},
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
		Time: time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("write messages: %w", err)
	}

	log.Printf("Kafka producer published transaction: %s", event.TransactionID)
	return nil
}

func (p *Producer) Close() error {
	log.Println("Closing Kafka producer")
	return p.writer.Close()
}
