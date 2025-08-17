package kafka

import (
	"context"
	"log/slog"
	"time"

	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaBatchTimeout = 50 * time.Millisecond
	kafkaMaxAttempts  = 10
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(addresses []string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(addresses...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: kafkaBatchTimeout,
		// At-least-once delivery
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		MaxAttempts:  kafkaMaxAttempts,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Produce(message, topic, key string) error {
	const mark = "Clients.Kafka.Producer"

	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(message),
			Time:  time.Now(),
		},
	)

	if err != nil {
		slog.With(slog.String("mark", mark)).Error("failed to produce event", slogext.Err(err))
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
