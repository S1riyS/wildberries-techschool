package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IKafkaHandler interface {
	// HandleMessage processes a Kafka message.
	// - ctx: Context for managing lifecycle and cancellation.
	// - kafkaMsg: The message received from Kafka as a byte slice.
	// - offset: The Kafka offset of the message.
	// - consumerNumber: The identifier for the consumer processing the message.
	// Returns an error if processing fails.
	HandleMessage(ctx context.Context, kafkaMsg []byte, offset kafka.Offset, consumerNumber int) error
}

// Constants for configuring Kafka consumer behavior.
const (
	sessionTimeout = 7000 // Session timeout in milliseconds. If no heartbeat is received within this time, the consumer will be considered inactive.
	noTimeout      = -1   // Indicates no timeout for reading messages from Kafka.
)

// Consumer represents a Kafka consumer that reads messages from a topic and processes them using a handler.
type Consumer struct {
	consumer       *kafka.Consumer // The underlying Kafka consumer instance.
	stop           bool            // Indicates whether the consumer should stop processing messages.
	handler        IKafkaHandler   // The handler for processing Kafka messages.
	consumerNumber int             // Identifier for this consumer instance (useful in multi-consumer scenarios).
}

func NewConsumer(handler IKafkaHandler, addresses []string, topic string, consumerGroup string, consumerNumber int) (*Consumer, error) {
	// Kafka consumer configuration.
	config := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(addresses, ","), // List of Kafka broker addresses.
		"group.id":                 consumerGroup,                // Consumer group ID for managing offsets and load balancing.
		"session.timeout.ms":       sessionTimeout,               // Timeout for detecting inactive consumers.
		"enable.auto.offset.store": false,                        // Prevent automatic offset storage; manual storage is used instead.
		"enable.auto.commit":       true,                         // Automatically commit offsets at intervals.
		"auto.commit.interval.ms":  5000,                         // Interval for automatic offset commits.
		"auto.offset.reset":        "earliest",                   // Reset behavior for new consumers (start from the earliest offset).
	}

	// Create a new Kafka consumer.
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Subscribe the consumer to the specified topic.
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return &Consumer{
		consumer:       consumer,
		handler:        handler,
		consumerNumber: consumerNumber,
	}, nil
}

// Start begins consuming messages from Kafka and processing them using the provided handler.
// - `ctx`: The context for controlling the lifecycle of the consumer.
func (c *Consumer) Start(ctx context.Context) {
	const mark = "infrastructure.kafka.Consumer.Start"

	logger := slog.With(slog.String("mark", mark))

	for {
		if c.stop {
			break // Stop consuming messages if the consumer is stopped.
		}
		// Read a message from Kafka with no timeout.
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logger.Warn("Failed to read message", slogext.Err(err))
			continue
		}
		if kafkaMsg == nil {
			continue // Skip processing if the message is nil.
		}

		// Handle the Kafka message using the provided handler.
		if err = c.handler.HandleMessage(ctx, kafkaMsg.Value, kafkaMsg.TopicPartition.Offset, c.consumerNumber); err != nil {
			logger.Warn("failed to handle message", slogext.Err(err))
		}

		// Store the message's offset to ensure it can be committed later.
		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			logger.Warn("Failed to store message", slogext.Err(err))
		}
	}
}

// Stop gracefully stops the Kafka consumer and commits the latest offsets.
// Returns an error if the commit or consumer close operation fails.
func (c *Consumer) Stop() error {
	c.stop = true // Signal the consumer to stop processing.
	if _, err := c.consumer.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return c.consumer.Close() // Close the consumer connection.
}
