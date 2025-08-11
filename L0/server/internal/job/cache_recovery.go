package job

import (
	"context"
	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
)

const ordersToRecover = 1000

func RecoverCache(ctx context.Context, repo domain.IOrderRepository, cache domain.IOrderCache) error {
	const mark = "job.RecoverCache"
	logger := slog.With(slog.String("mark", mark))

	orders, err := repo.GetRecentlyCreated(ctx, ordersToRecover)
	if err != nil {
		logger.Error("Failed to get recently created orders from repository", slogext.Err(err))
		return err
	}

	counter := 0
	for _, order := range orders {
		err = cache.Save(ctx, order)
		if err != nil {
			logger.Error("Failed to save order to cache", slog.String("order_uid", order.OrderUID), slogext.Err(err))
			return err
		}
		counter++
	}

	logger.Info("Recovered cache", slog.Int("orders_count", counter))
	return nil
}
