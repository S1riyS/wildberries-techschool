package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/app"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogpretty"
	"github.com/joho/godotenv"
)

const (
	localEnvPath = "configs/.env.local"
)

func main() {
	const mark = "server.main"

	// Read flags
	isDevMode := flag.Bool("local", false, "Local mode")
	flag.Parse()

	// Load local .env file in dev mode
	if *isDevMode {
		err := godotenv.Load(localEnvPath)
		if err != nil {
			panic(fmt.Errorf("error loading .env file: %w", err))
		}
	}

	// Init config
	cfg := config.MustNew()

	// Init and set logger
	slog.SetDefault(setupLogger(cfg.Env))

	// DEBUG: Print config
	if cfg.Env != config.EnvProd {
		slog.With(slog.String("mark", mark)).Debug("Config", slog.Any("config", cfg))
	}

	// Init and run app
	ctx := context.Background()
	application := app.MustNew(cfg)
	go application.MustRun(ctx) // Run app in separate goroutine

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	slog.With(slog.String("mark", mark)).Info("Gracefully stopped")
}

func setupLogger(env config.EnvType) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = setupPrettySlog()
	// TODO: change back to production logger
	// case config.EnvProd:
	// 	log = slog.New(
	// 		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	// 	)
	// }
	case config.EnvProd:
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
