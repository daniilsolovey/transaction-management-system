package kafka_broker

import (
	"log/slog"
	"strings"

	"github.com/daniilsolovey/transaction-management-system/internal/usecase"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	topic   = "transcations"
	groupID = "transaction-group"
)

type Consumer struct {
	reader  *kafka.Reader
	usecase *usecase.TransactionUseCase
	log     *slog.Logger
}

type Writer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

type Broker struct {
	Writer   *Writer
	Consumer *Consumer
}

func NewBroker(useCase *usecase.TransactionUseCase, logger *slog.Logger) *Broker {
	brokers := strings.Split(viper.GetString("KAFKA_ADDRESS"), ",")

	writer := newWriter(brokers, topic, logger)
	consumer := newConsumer(brokers, topic, groupID, useCase, logger)

	return &Broker{
		Writer:   writer,
		Consumer: consumer,
	}
}

func (b *Broker) Close() error {
	if err := b.Writer.Close(); err != nil {
		return err
	}
	return b.Consumer.reader.Close()
}

func (w *Writer) Close() error {
	return w.writer.Close()
}

func newWriter(brokers []string, topic string, log *slog.Logger) *Writer {
	return &Writer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		log: log,
	}
}

func newConsumer(brokers []string, topic, groupID string,
	uc *usecase.TransactionUseCase, log *slog.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	return &Consumer{reader: r, usecase: uc, log: log}
}
