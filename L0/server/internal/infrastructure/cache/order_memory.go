package cache

import (
	"fmt"
	"sync"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
)

type OrderInMemoryCache struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func NewOrderInMemoryCache() *OrderInMemoryCache {
	return &OrderInMemoryCache{
		orders: make(map[string]*domain.Order),
	}
}

func (c *OrderInMemoryCache) Save(order *domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
	return nil
}

func (c *OrderInMemoryCache) Get(orderID string) (*domain.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}
