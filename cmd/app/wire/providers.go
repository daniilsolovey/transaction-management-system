package wire

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/daniilsolovey/transaction-management-system/internal/pkg/grpcserver"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/postgres"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func ProvidePostgres(logger *slog.Logger) (*postgres.Repository, error) {
	ctx := context.Background()
	url := viper.GetString("DATABASE_URL")

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = int32(viper.GetInt("DB_MAX_CONNS"))
	lifetime, _ := time.ParseDuration(viper.GetString("DB_MAX_CONN_LIFETIME"))
	cfg.MaxConnLifetime = lifetime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return postgres.New(pool, logger), nil
}

func ProvideRedis(logger *slog.Logger) (*redis.Repository, error) {
	host := viper.GetString("REDIS_HOST")    // "redis"
	port := viper.GetString("REDIS_PORT")    // "6379"
	addr := fmt.Sprintf("%s:%s", host, port) // "redis:6379"

	password := viper.GetString("REDIS_PASSWORD")
	db := viper.GetInt("REDIS_DB")

	return redis.New(addr, password, db, logger), nil
}

func ProvideRPCServer() (*grpcserver.Server, error) {
	return grpcserver.New(viper.GetInt("GRPC_SERVER_PORT")), nil
}

func ProvideLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo},
		),
	)
}
