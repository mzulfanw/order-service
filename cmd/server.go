package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/order-service/configs"
	"github.com/mzulfanw/order-service/internal/handler"
	"github.com/mzulfanw/order-service/internal/repository"
	"github.com/mzulfanw/order-service/internal/service"
	"github.com/mzulfanw/order-service/pkg"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.InfoLevel)

	loadConfig, err := configs.LoadConfig()
	if err != nil {
		log.Infof("Failed to load config: %v", err)
	}

	database, err := pkg.InitDatabase(loadConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	cache := pkg.NewRedisCache(loadConfig.RedisConfig)
	mq := pkg.NewRabbitMQ(loadConfig.RabbitMQConfig)
	client := pkg.NewProductClient(loadConfig.ProductServiceUrl)

	orderRepository := repository.NewOrderRepository(database)
	orderService := service.NewOrderService(orderRepository, cache, client, mq)

	orderHandler := handler.NewOrderHandler(orderService)

	r := mux.NewRouter()
	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/product/{productId}", orderHandler.GetByProductID).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", loadConfig.Port),
		Handler: r,
	}

	go func() {
		log.Infof("🚀 order-service running on port %d", 3001)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("server forced to shutdown: %v", err)
	}

	mq.Close()
	_ = database.Close()
	log.Info("✅ order-service stopped cleanly")
}
