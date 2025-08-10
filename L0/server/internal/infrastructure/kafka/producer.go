package kafka

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// unknownType represents an error for unrecognized event types.
var errUnknownType = errors.New("unknown event type")

const (
	// flushTimeout defines the timeout in milliseconds for flushing the Kafka producer.
	flushTimeout = 5000
)

// Producer wraps a Kafka producer to handle message production to Kafka topics.
type Producer struct {
	// producer is the underlying Kafka producer instance.
	producer *kafka.Producer
}

// NewProducer creates and initializes a new Producer with the given Kafka broker addresses.
func NewProducer(addresses []string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addresses, ","),
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	return &Producer{producer: producer}, nil
}

// Produce sends a message with a key and timestamp to the specified Kafka topic.
func (p *Producer) Produce(message, topic, key string) error {
	const mark = "Cients.Kafka.Producer"

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
		Key:            []byte(key),
		Timestamp:      time.Now(),
	}
	kafkaChan := make(chan kafka.Event)

	if err := p.producer.Produce(kafkaMessage, kafkaChan); err != nil {
		slog.With(slog.String("mark", mark)).Error("failed to produce event", slogext.Err(err))
		return err
	}

	event := <-kafkaChan
	switch event.(type) {
	case *kafka.Message:
		m := event.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
		return nil
	default:
		return errUnknownType
	}
}

// Close gracefully shuts down the Kafka producer after flushing pending messages.
func (p *Producer) Close() error {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
	return nil
}
