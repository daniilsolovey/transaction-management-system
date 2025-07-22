package kafka_broker

import (
	"context"
	"encoding/json"

	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/segmentio/kafka-go"
)

// RunConsumer starts the consumer loop.
func (c *Consumer) RunConsumer(ctx context.Context) {
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
				c.reader.CommitMessages(ctx, msg)
				continue
			}

			err = c.usecase.Create(ctx, message)
			if err != nil {
				c.log.Error("failed to process transaction", "error", err)
				// Retry later — don't commit
				continue
			}

			err = c.reader.CommitMessages(ctx, msg)
			if err != nil {
				c.log.Error("failed to commit message", "err", err)
			}
		}
	}
}

func (w *Writer) Publish(ctx context.Context,
	message domain.CreateTransactionMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		w.log.Error("failed to marshal Kafka message", "err", err)
		return err
	}

	err = w.writer.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		w.log.Error("failed to publish to Kafka", "err", err)
		return err
	}

	w.log.Info("message published to Kafka", "user_id", message.UserID)
	return nil
}
