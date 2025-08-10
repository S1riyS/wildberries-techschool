package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/kafka"
	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
)

func main() {
	// Read flags
	var envPath string
	flag.StringVar(&envPath, "env-path", "", "Path to .env file")
	flag.Parse()

	// Load local .env file in dev mode
	err := godotenv.Load(envPath)
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}
	log.Println("Loaded .env file")

	// Init config
	cfg := config.MustNew()
	log.Println("Loaded config")

	// Create order producer
	orderProducer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		panic(fmt.Errorf("error creating order producer: %v", err))
	}
	log.Println("Created order producer")
	defer orderProducer.Close()

	// Generate order
	var randomOrder domain.Order
	err = faker.FakeData(&randomOrder)
	if err != nil {
		panic(fmt.Errorf("error generating order: %v", err))
	}
	orderJSON, err := json.Marshal(randomOrder)
	if err != nil {
		panic(fmt.Errorf("error marshaling order to json: %v", err))
	}
	log.Println("Created order JSON")

	err = orderProducer.Produce(string(orderJSON), cfg.Kafka.Topic, "")
	if err != nil {
		panic(fmt.Errorf("error producing order: %v", err))
	}
}
