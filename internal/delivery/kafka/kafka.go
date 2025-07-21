package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/segmentio/kafka-go"
)

type TransactionUseCase interface {
	Create(ctx context.Context, dto domain.CreateTransactionMessage) error
}

type Consumer struct {
	reader  *kafka.Reader
	usecase TransactionUseCase
	log     *slog.Logger
}

func NewConsumer(brokers []string, topic, groupID string, uc TransactionUseCase, log *slog.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	return &Consumer{reader: r, usecase: uc, log: log}
}

// Run starts the consumer loop.
func (c *Consumer) Run(ctx context.Context) {
	c.log.Info("Starting Kafka consumer", "topic", c.reader.Config().Topic)
	defer c.reader.Close()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Stopping Kafka consumer...")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.log.Error("could not fetch message", "error", err)
				continue
			}

			var message domain.CreateTransactionMessage
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				c.log.Error("failed to unmarshal message", "error", err)
				c.reader.CommitMessages(ctx, msg) // Commit poison pill to avoid reprocessing
				continue
			}

			if err := c.usecase.Create(ctx, message); err != nil {
				c.log.Error("failed to process transaction", "error", err)
				// Here you would implement a retry or dead-letter queue strategy
			}

			c.reader.CommitMessages(ctx, msg)
		}
	}
}
