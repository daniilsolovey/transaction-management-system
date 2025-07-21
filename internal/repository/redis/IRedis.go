package redis

import (
	"context"
	"log/slog"

	"github.com/go-redis/redis"
)

type Repository struct {
	client *redis.Client
	log    *slog.Logger
}

type IRepository interface {
	GetValueByKey(ctx context.Context, key string) (int, error)
	Del(ctx context.Context, key string) error
	Close() error
}

func New(addr, password string, db int, log *slog.Logger) *Repository {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &Repository{client: rdb, log: log}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.client.Ping().Err()
}

func (r *Repository) Close() error {
	return r.client.Close()
}
