package app

import (
	"context"

	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
)

type App struct {
	logger     *slog.Logger
	config     config.Config
	httpServer *HTTPServer
}

func New(_ context.Context, logger *slog.Logger, cfg config.Config) *App {
	// const mark = "app.New"

	app := &App{
		logger:     logger,
		config:     cfg,
		httpServer: NewHTTPServer(logger, cfg),
	}

	app.initValidator()

	return app
}

func (a *App) MustRun() {
	const mark = "app.Run"

	logger := a.logger.With(slog.String("mark", mark))

	if err := a.httpServer.Run(); err != nil {
		logger.Error("failed to start http server", slog.Int("port", a.config.HTTP.Port), slogext.Err(err))
		panic(err)
	}
}

func (a *App) Stop() {
	a.httpServer.Stop()
}

func (a *App) initValidator() {
	const mark = "app.initValidator"

	a.logger.Warn("Validator is NOT initialized", slog.String("mark", mark))
}
