package service

import (
	"context"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
)

type OrderService struct {
	repo domain.IOrderRepository
}

func NewOrderService(repo domain.IOrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Save(ctx context.Context, order *domain.Order) error {
	return s.repo.Save(ctx, order)
}

func (s *OrderService) GetOne(ctx context.Context, orderID string) (*domain.Order, error) {
	return s.repo.Get(ctx, orderID)
}
