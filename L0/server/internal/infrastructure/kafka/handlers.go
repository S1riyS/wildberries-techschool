package kafka

import (
	"context"
	"encoding/json"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) HandleMessage(ctx context.Context, kafkaMsg []byte, offset kafka.Offset, consumerNumber int) error {
	// Unmarshal the Kafka message into an Order struct.
	var order domain.Order
	if err := json.Unmarshal(kafkaMsg, &order); err != nil {
		return err
	}

	return h.service.Save(ctx, &order)
}
