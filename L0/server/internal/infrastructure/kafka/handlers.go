package kafka

import (
	"context"
	"encoding/json"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) HandleMessage(ctx context.Context, kafkaMsg []byte, offset int64, consumerNumber int) error {
	var order domain.Order
	if err := json.Unmarshal(kafkaMsg, &order); err != nil {
		return err
	}

	return h.service.Save(ctx, &order)
}