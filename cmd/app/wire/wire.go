//go:build wireinject
// +build wireinject

package wire

import (
	"log/slog"

	"github.com/daniilsolovey/transaction-management-system/internal/repository/postgres"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/redis"
	"github.com/google/wire"
)

type Service struct {
	Postgres *postgres.Repository
	Redis    *redis.Repository
	Logger   *slog.Logger
}

func Initialize() (*Service, func(), error) {
	wire.Build(
		ProvideLogger,
		ProvidePostgres,
		ProvideRedis,
		wire.Struct(new(Service), "*"),
	)
	return nil, nil, nil
}
