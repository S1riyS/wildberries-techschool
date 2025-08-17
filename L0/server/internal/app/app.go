package app

import (
	"context"
	"sync"

	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/job"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
)

type App struct {
	config   config.Config
	resolver *Resolver
}

func MustNew(cfg config.Config) *App {
	return &App{
		config:   cfg,
		resolver: NewResolver(cfg),
	}
}

func (a *App) MustRun(ctx context.Context) {
	const mark = "app.Run"

	logger := slog.With(slog.String("mark", mark))

	const numTasks = 3
	wg := &sync.WaitGroup{}
	wg.Add(numTasks)

	// Recover cache
	go func() {
		defer wg.Done()
		err := job.RecoverCache(ctx, a.resolver.OrderRepository(ctx), a.resolver.OrderCache())
		if err != nil {
			logger.Error("failed to recover cache", slogext.Err(err))
		}
	}()

	// Run HTTP server
	go func() {
		defer wg.Done()
		err := a.resolver.HTTPServer(ctx).Run()
		if err != nil {
			logger.Error("failed to start http server", slog.Int("port", a.config.HTTP.Port), slogext.Err(err))
		}
	}()

	// Run kafka consumers
	go func() {
		defer wg.Done()
		a.resolver.OrderConsumer(ctx).Start(ctx)
	}()

	wg.Wait()
}

func (a *App) Stop() {
	ctx := context.Background()
	// Stop HTTP server
	a.resolver.HTTPServer(ctx).Stop()

	// Stop kafka consumer
	err := a.resolver.OrderConsumer(ctx).Stop()
	if err != nil {
		slog.Error("failed to stop kafka consumer", slogext.Err(err))
	}
}
