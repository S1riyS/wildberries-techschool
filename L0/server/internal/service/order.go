package service

import "github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"

type OrderService struct {
	repo domain.IOrderRepository
}

func NewOrderService(repo domain.IOrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Save(order *domain.Order) error {
	return s.repo.Save(order)
}

func (s *OrderService) Get(orderID string) (*domain.Order, error) {
	return s.repo.Get(orderID)
}
