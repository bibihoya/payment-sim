package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	topic  string
	group  string
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	log.Printf("Kafka consumer created for topic: %s, group: %s", topic, groupID)
	return &Consumer{
		reader: reader,
		topic:  topic,
		group:  groupID,
	}
}

func (c *Consumer) Consume(ctx context.Context,
	handler func(context.Context, *TransactionEvent) error) error {
	log.Printf("Consumer for topic: %s, group: %s started", c.topic, c.group)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer closed", c.topic, c.group)
			return ctx.Err()

		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			event, err := TransEventFromJSON(msg.Value)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			log.Printf("Kafka producer published transaction: %s", event.TransactionID)
			err = handler(ctx, event)
			if err != nil {
				log.Printf("Error processing transaction: %v", err)
				continue
			}
		}
	}
}

func (c *Consumer) Close() error {
	log.Println("Closing Kafka consumer")
	return c.reader.Close()
}
