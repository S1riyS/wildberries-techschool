package app

import (
	"context"

	v1 "github.com/S1riyS/wildberries-techschool/L0/server/internal/api/http/handler/v1"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/domain"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/cache"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/kafka"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/storage"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
)

const (
	orderCacheCapacity = 1000 // TODO: move to configuration (?)
)

type Resolver struct {
	config config.Config

	dbClient   postgresql.Client
	httpServer *HTTPServer

	orderCache      domain.IOrderCache
	orderRepository domain.IOrderRepository
	orderService    *service.OrderService
	orderHandler    *v1.OrderHandler
	orderConsumer   *kafka.Consumer
}

func NewResolver(cfg config.Config) *Resolver {
	return &Resolver{
		config: cfg,
	}
}

func (r *Resolver) DBClient(ctx context.Context) postgresql.Client {
	if r.dbClient == nil {
		r.dbClient = postgresql.MustNewClient(ctx, r.config.Database)
	}
	return r.dbClient
}

func (r *Resolver) HTTPServer(ctx context.Context) *HTTPServer {
	if r.httpServer == nil {
		r.httpServer = NewHTTPServer(r.config, r.OrderHandler(ctx))
	}
	return r.httpServer
}

func (r *Resolver) OrderCache() domain.IOrderCache {
	if r.orderCache == nil {
		r.orderCache = cache.NewOrderInMemoryCache(orderCacheCapacity)
	}
	return r.orderCache
}

func (r *Resolver) OrderRepository(ctx context.Context) domain.IOrderRepository {
	if r.orderRepository == nil {
		r.orderRepository = storage.NewOrderRepository(r.DBClient(ctx), r.OrderCache())
	}
	return r.orderRepository
}

func (r *Resolver) OrderService(ctx context.Context) *service.OrderService {
	if r.orderService == nil {
		r.orderService = service.NewOrderService(r.OrderRepository(ctx))
	}
	return r.orderService
}

func (r *Resolver) OrderHandler(ctx context.Context) *v1.OrderHandler {
	if r.orderHandler == nil {
		r.orderHandler = v1.NewOrderHandler(r.OrderService(ctx))
	}
	return r.orderHandler
}

func (r *Resolver) OrderConsumer(ctx context.Context) *kafka.Consumer {
	if r.orderConsumer == nil {
		consumer, err := kafka.NewConsumer(
			kafka.NewOrderHandler(r.OrderService(ctx)),
			r.config.Kafka.Brokers,
			r.config.Kafka.Topic,
			"order",
		)
		if err != nil {
			panic(err)
		}
		r.orderConsumer = consumer
	}
	return r.orderConsumer
}
