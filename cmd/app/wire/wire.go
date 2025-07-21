//go:build wireinject
// +build wireinject

package wire

import (
	"log/slog"

	"github.com/google/wire"
)

type Service struct {
	Postgres  *postgres.Repository
	Redis     *redis.Repository
	RPCServer *grpcserver.Server
	Logger    *slog.Logger
}

func Initialize() (*Service, func(), error) {
	wire.Build(
		ProvideLogger,
		ProvidePostgres,
		ProvideRedis,
		ProvideRPCServer,
		wire.Struct(new(Service), "*"),
	)
	return nil, nil, nil
}
