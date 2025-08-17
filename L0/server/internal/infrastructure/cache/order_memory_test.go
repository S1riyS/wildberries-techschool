package cache

import (
	"context"
	"sync"
	"testing"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderInMemoryCache(t *testing.T) {
	ctx := context.Background()
	cache := NewOrderInMemoryCache(2)

	order1 := &domain.Order{OrderUID: "order1"}
	order2 := &domain.Order{OrderUID: "order2"}
	order3 := &domain.Order{OrderUID: "order3"}

	t.Run("empty cache", func(t *testing.T) {
		_, err := cache.Get(ctx, "nonexistent")
		assert.Error(t, err, "should return error for non-existent order")
	})

	t.Run("save and get", func(t *testing.T) {
		err := cache.Save(ctx, order1)
		require.NoError(t, err)

		got, err := cache.Get(ctx, order1.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, order1, got, "should return saved order")
	})

	t.Run("update existing", func(t *testing.T) {
		updatedOrder := &domain.Order{OrderUID: "order1", TrackNumber: "updated"}
		err := cache.Save(ctx, updatedOrder)
		require.NoError(t, err)

		got, err := cache.Get(ctx, order1.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, updatedOrder, got, "should update existing order")
	})

	t.Run("lru eviction", func(t *testing.T) {
		// Fill cache to capacity
		err := cache.Save(ctx, order1)
		require.NoError(t, err)
		err = cache.Save(ctx, order2)
		require.NoError(t, err)

		// Access order1 to make it recently used
		_, err = cache.Get(ctx, order1.OrderUID)
		require.NoError(t, err)

		// Add order3 - should evict order2 (least recently used)
		err = cache.Save(ctx, order3)
		require.NoError(t, err)

		// order2 should be evicted
		_, err = cache.Get(ctx, order2.OrderUID)
		assert.Error(t, err, "order2 should be evicted")

		// order1 and order3 should still be in cache
		_, err = cache.Get(ctx, order1.OrderUID)
		assert.NoError(t, err)
		_, err = cache.Get(ctx, order3.OrderUID)
		assert.NoError(t, err)
	})

	t.Run("concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		for range 100 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = cache.Save(ctx, order1)
				_, _ = cache.Get(ctx, order1.OrderUID)
			}()
		}
		wg.Wait()

		got, err := cache.Get(ctx, order1.OrderUID)
		require.NoError(t, err)
		assert.Equal(t, order1, got, "should handle concurrent access")
	})
}
