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
)

func init() {
	configs.Init()
}

func main() {
	service, cleanup, err := wire.Initialize()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	ctx := context.Background()

	// Check connection PostgreSQL
	if err := service.Postgres.Ping(ctx); err != nil {
		slog.Error("PostgreSQL not available", "error", err)
		os.Exit(1)
	}

	// Check connection Redis
	if err := service.Redis.Ping(ctx); err != nil {
		slog.Error("Redis not available", "error", err)
		os.Exit(1)
	}

	repo := repository.New(service.Postgres, service.Redis)
	useCase := usecase.NewTransactionUseCase(repo, service.Logger)
	handler := deliveryhttp.NewTransactionHandler(useCase, service.Logger)

	engine := handler.RegisterRoutes()
	port := viper.GetInt("HTTP_PORT")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run HTTP-server
	go func() {
		slog.Info("HTTP server started", "port", port)
		if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil &&
			err != http.ErrServerClosed {
			slog.Error("HTTP server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("service stopping")
}
