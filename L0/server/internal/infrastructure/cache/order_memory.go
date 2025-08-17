package cache

import (
	"container/list"
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
)

// OrderInMemoryCache is an in-memory LRU cache for orders. Threadsafe due to sync.RWMutex.
type OrderInMemoryCache struct {
	mu       sync.RWMutex
	orders   map[string]*list.Element // Maps order UID to list element
	list     *list.List               // Maintains order of access (front = most recent)
	capacity int                      // Maximum number of orders to cache
}

func NewOrderInMemoryCache(capacity int) *OrderInMemoryCache {
	return &OrderInMemoryCache{
		orders:   make(map[string]*list.Element),
		list:     list.New(),
		capacity: capacity,
	}
}

// cacheEntry represents an item in the cache.
type cacheEntry struct {
	key   string // Order UID
	order *domain.Order
}

func (c *OrderInMemoryCache) Save(_ context.Context, order *domain.Order) error {
	const mark = "OrderInMemoryCache.Save"
	logger := slog.With(slog.String("mark", mark))

	c.mu.Lock()
	defer c.mu.Unlock()

	// If order already exists, update it and move to front
	if elem, exists := c.orders[order.OrderUID]; exists {
		c.list.MoveToFront(elem)
		cacheEntry, ok := elem.Value.(*cacheEntry)
		if !ok {
			logger.Error("Failed to cast cache entry", slog.String("order_uid", order.OrderUID))
			return errors.New("failed to cast cache entry")
		}
		cacheEntry.order = order
		logger.Debug("Order updated in cache", slog.String("order_uid", order.OrderUID))
		return nil
	}

	// If at capacity, remove the least recently used order
	if len(c.orders) >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			oldestEntry, ok := oldest.Value.(*cacheEntry)
			if !ok {
				logger.Error("Failed to cast cache entry")
				return errors.New("failed to cast cache entry")
			}

			delete(c.orders, oldestEntry.key)
			c.list.Remove(oldest)
			logger.Debug("Evicted LRU order from cache", slog.String("order_uid", oldestEntry.key))
		}
	}

	// Add new order to cache
	entry := &cacheEntry{
		key:   order.OrderUID,
		order: order,
	}
	elem := c.list.PushFront(entry)
	c.orders[order.OrderUID] = elem

	logger.Debug("Order saved to cache", slog.String("order_uid", order.OrderUID))
	return nil
}

func (c *OrderInMemoryCache) Get(_ context.Context, orderID string) (*domain.Order, error) {
	const mark = "OrderInMemoryCache.Get"
	logger := slog.With(slog.String("mark", mark))

	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.orders[orderID]
	if !exists {
		logger.Debug("Order not found in cache", slog.String("order_uid", orderID))
		return nil, errors.New("order not found")
	}

	// Move the accessed order to front (most recently used)
	c.list.MoveToFront(elem)
	entry, ok := elem.Value.(*cacheEntry)
	if !ok {
		logger.Error("Failed to cast cache entry", slog.String("order_uid", orderID))
		return nil, errors.New("failed to cast cache entry")
	}
	order := entry.order

	logger.Debug("Order found in cache", slog.String("order_uid", order.OrderUID))
	return order, nil
}
