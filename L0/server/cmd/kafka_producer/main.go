package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/kafka"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogpretty"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/joho/godotenv"
)

const (
	localEnvPath = "configs/.env.local"
)

func main() {
	const mark = "kafka_producer.main"

	logger := setupPrettySlog()
	logger = logger.With(slog.String("mark", mark))

	// Read flags
	isDevMode := flag.Bool("local", false, "Local mode")
	flag.Parse()

	// Load local .env file in dev mode
	if *isDevMode {
		err := godotenv.Overload(localEnvPath)
		if err != nil {
			logger.Error("Error loading .env file", slogext.Err(err))
			return
		}
	}

	// Init config
	cfg := config.MustNew()
	logger.Info("Created config",
		slog.String("kafkaBrokers", fmt.Sprintf("%v", cfg.Kafka.Brokers)),
	)

	// Create order producer
	orderProducer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		logger.Error("Error creating order producer", slogext.Err(err))
		return
	}
	defer orderProducer.Close()

	// Generate order
	var randomOrder domain.Order
	const minSize = 1
	const maxSize = 3
	err = faker.FakeData(&randomOrder, options.WithRandomMapAndSliceMinSize(minSize), options.WithRandomMapAndSliceMaxSize(maxSize))
	if err != nil {
		logger.Error("Error generating random order", slogext.Err(err))
		return
	}
	orderJSON, err := json.Marshal(randomOrder)
	if err != nil {
		logger.Error("Error marshaling order to json", slogext.Err(err))
		return
	}

	err = orderProducer.Produce(string(orderJSON), cfg.Kafka.Topic, randomOrder.OrderUID)
	if err != nil {
		logger.Error("Error producing order", slogext.Err(err))
	}

	logger.Info("Produced order", slog.String("orderUID", randomOrder.OrderUID))
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
