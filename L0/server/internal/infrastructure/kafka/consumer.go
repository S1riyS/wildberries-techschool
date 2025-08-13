package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/segmentio/kafka-go"
)

type IKafkaHandler interface {
	HandleMessage(ctx context.Context, kafkaMsg []byte, offset int64, consumerNumber int) error
}

type Consumer struct {
	reader         *kafka.Reader
	stop           bool
	handler        IKafkaHandler
	consumerNumber int
}

func NewConsumer(handler IKafkaHandler, addresses []string, topic string, consumerGroup string) (*Consumer, error) {
	const mark = "infrastructure.kafka.NewConsumer"

	logger := slog.With(slog.String("mark", mark))
	logger.Info("Creating kafka consumer",
		slog.String("brokers", fmt.Sprintf("%v", addresses)),
		slog.String("topic", topic),
		slog.String("group", consumerGroup),
	)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addresses,
		GroupID:        consumerGroup,
		Topic:          topic,
		SessionTimeout: 7 * time.Second,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0, // Disable auto-commit
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	const mark = "infrastructure.kafka.Consumer.Start"
	logger := slog.With(slog.String("mark", mark))

	logger.Info("Starting kafka consumer")
	for {
		if c.stop {
			break
		}

		// Read message
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logger.Warn("Failed to read message", slogext.Err(err))
			continue
		}

		// Handle message
		if err := c.handler.HandleMessage(ctx, msg.Value, msg.Offset, c.consumerNumber); err != nil {
			logger.Warn("Failed to handle message", slogext.Err(err))
		}

		// Commit message manually
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			// Обработка ошибки коммита
			logger.Warn("Failed to commit message", slogext.Err(err))
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	return c.reader.Close()
}
