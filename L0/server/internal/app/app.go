package app

import (
	"context"
	"sync"

	"log/slog"

	v1 "github.com/S1riyS/wildberries-techschool/L0/server/internal/api/http/handler/v1"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/cache"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/kafka"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/storage"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
)

type App struct {
	config        config.Config
	httpServer    *HTTPServer
	kafkaConsumer *kafka.Consumer
}

func MustNew(_ context.Context, cfg config.Config) *App {
	const mark = "app.New"

	// Dependencies
	dbClient := postgresql.MustNewClient(context.TODO(), cfg.Database)
	orderCache := cache.NewOrderInMemoryCache()
	orderRepository := storage.NewOrderRepository(dbClient, orderCache)
	orderService := service.NewOrderService(orderRepository)
	orderHandler := v1.NewOrderHandler(orderService)

	// Kafka consumer
	consumer, err := kafka.NewConsumer(
		kafka.NewOrderHandler(orderService),
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		"order",
	)
	if err != nil {
		slog.With(slog.String("mark", mark)).Error("failed to create kafka consumer", slogext.Err(err))
		panic(err)
	}

	app := &App{
		config:        cfg,
		httpServer:    NewHTTPServer(cfg, orderHandler),
		kafkaConsumer: consumer,
	}

	return app
}

func (a *App) MustRun() {
	const mark = "app.Run"

	logger := slog.With(slog.String("mark", mark))

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.httpServer.Run()
		if err != nil {
			logger.Error("failed to start http server", slog.Int("port", a.config.HTTP.Port), slogext.Err(err))
		}
	}()

	go func() {
		defer wg.Done()
		a.kafkaConsumer.Start(context.TODO())
	}()

	wg.Wait()
}

func (a *App) Stop() {
	a.httpServer.Stop()
	a.kafkaConsumer.Stop()
}
