package app

import (
	"context"

	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
)

type App struct {
	config     config.Config
	httpServer *HTTPServer
}

func New(_ context.Context, cfg config.Config) *App {
	// const mark = "app.New"

	app := &App{
		config:     cfg,
		httpServer: NewHTTPServer(cfg),
	}

	app.initValidator()

	return app
}

func (a *App) MustRun() {
	const mark = "app.Run"

	logger := slog.With(slog.String("mark", mark))

	if err := a.httpServer.Run(); err != nil {
		logger.Error("failed to start http server", slog.Int("port", a.config.HTTP.Port), slogext.Err(err))
		panic(err)
	}

	// TODO: move consumer creation somwhere else
	// consumer := kafka.NewConsumer(kafka.NewOrderHandler(service), a.config.Kafka.Brokers, a.config.Kafka.Topic, a.config.Kafka.ConsumerGroup, a.config.Kafka.ConsumerNumber)
	// consumer.Start(context.Background())
}

func (a *App) Stop() {
	a.httpServer.Stop()
}

func (a *App) initValidator() {
	const mark = "app.initValidator"

	slog.Warn("Validator is NOT initialized", slog.String("mark", mark))
}
