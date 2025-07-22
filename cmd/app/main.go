package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/daniilsolovey/transaction-management-system/cmd/app/wire"
	"github.com/daniilsolovey/transaction-management-system/configs"
	deliveryhttp "github.com/daniilsolovey/transaction-management-system/internal/delivery/http"

	"github.com/daniilsolovey/transaction-management-system/internal/repository"
	"github.com/daniilsolovey/transaction-management-system/internal/usecase"
	"github.com/spf13/viper"

	kafka_broker "github.com/daniilsolovey/transaction-management-system/internal/delivery/kafka"

	_ "github.com/daniilsolovey/transaction-management-system/docs"
)

func init() {
	configs.Init()
}

// @title Transaction Manager API
// @version 1.0
// @description API for managing user transactions via Kafka and PostgreSQL
// @host localhost:3000
// @BasePath /
// @schemes http

func main() {
	configs.Init()

	// Initialize Postgres, Redis, Logger using wire
	service, cleanup, err := wire.Initialize()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ping services
	if err := service.Postgres.Ping(ctx); err != nil {
		slog.Error("PostgreSQL not available", "error", err)
		os.Exit(1)
	}
	if err := service.Redis.Ping(ctx); err != nil {
		slog.Error("Redis not available", "error", err)
		os.Exit(1)
	}

	// Dependencies

	repo := repository.New(service.Postgres, service.Redis)
	useCase := usecase.NewTransactionUseCase(repo, service.Logger)

	// Initialize Kafka broker using package
	broker := kafka_broker.NewBroker(useCase, service.Logger)
	defer broker.Close()

	// Register HTTP routes using package
	handler := deliveryhttp.NewTransactionHandler(useCase, service.Logger, broker.Writer)
	engine := handler.RegisterRoutes()

	// Run HTTP server
	port := viper.GetInt("HTTP_PORT")
	go func() {
		slog.Info("HTTP server started", "port", port)
		if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "err", err)
			os.Exit(1)
		}
	}()

	// Run Kafka consumer
	go broker.Consumer.RunConsumer(ctx)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")
}
