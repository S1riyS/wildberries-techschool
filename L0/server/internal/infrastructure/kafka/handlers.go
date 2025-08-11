package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
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
	const mark = "kafka.OrderHandler.HandleMessage"
	logger := slog.With(slog.String("mark", mark))

	var order domain.Order
	if err := json.Unmarshal(kafkaMsg, &order); err != nil {
		logger.Error("Failed to unmarshal order", slogext.Err(err))
		return err
	}

	err := order.Validate()
	if err != nil {
		logger.Error("Failed to validate order", slogext.Err(err))
		return err
	}

	return h.service.Save(ctx, &order)
}
