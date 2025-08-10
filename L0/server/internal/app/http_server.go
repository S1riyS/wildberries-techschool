package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	v1 "github.com/S1riyS/wildberries-techschool/L0/server/internal/api/http/handler/v1"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/config"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/cache"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/infrastructure/storage"
	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/logger/slogext"
	"github.com/S1riyS/wildberries-techschool/L0/server/pkg/postgresql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	config config.Config

	ginInstance *gin.Engine // Gin engine that runs on `httpSrv`
	// httpSrv     *http.Server // Underlying HTTP server
}

func NewHTTPServer(config config.Config) *HTTPServer {
	server := &HTTPServer{
		config: config,
	}

	// Run init steps
	server.initGin()
	server.initControllers()

	return server
}

// Run starts the http server.
//
// It runs the gin.Engine and returns an error if it can't start the server.
// The port number is taken from the configuration.
func (hs *HTTPServer) Run() error {
	const mark = "httpServer.Run"

	logger := slog.With(slog.String("mark", mark))

	port := fmt.Sprintf(":%d", hs.config.HTTP.Port)
	if err := hs.ginInstance.Run(port); err != nil {
		logger.Error("failed to start http server", slog.Int("port", hs.config.HTTP.Port), slogext.Err(err))
	}

	return nil
}

func (hs *HTTPServer) Stop() {
	const mark = "httpServer.Stop"

	logger := slog.With(slog.String("mark", mark))
	logger.Warn("httpServer.Stop is NOT implemented yet", slog.Int("port", hs.config.HTTP.Port))
}

func (hs *HTTPServer) initGin() {
	hs.ginInstance = gin.New()

	if hs.config.Env == config.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS configuration
	hs.ginInstance.Use(cors.New(cors.Config{
		AllowOrigins:     hs.config.HTTP.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, //nolint:mnd // Later will be retrieved from config (probably)
	}))

	// Middlewares
	hs.ginInstance.Use(
		gin.Recovery(),
		gin.Logger(),
	)
}

func (hs *HTTPServer) initControllers() {
	// API
	apiGroup := hs.ginInstance.Group("/api")
	v1Group := apiGroup.Group("/v1")

	orderGroup := v1Group.Group("/orders")

	dbClient := postgresql.MustNewClient(context.TODO(), hs.config.Database)

	orderCache := cache.NewOrderInMemoryCache()
	orderRepository := storage.NewOrderRepository(dbClient, orderCache)
	orderService := service.NewOrderService(orderRepository)
	orderHandler := v1.NewOrderHandler(orderService)
	{
		orderGroup.GET("/:id", orderHandler.GetOne)
	}
}
