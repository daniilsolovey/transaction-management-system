package repository

import (
	"github.com/daniilsolovey/transaction-management-system/internal/repository/postgres"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/redis"
)

type repo struct {
	pg  postgres.IRepository
	red redis.IRepository
}

func New(pg postgres.IRepository, red redis.IRepository) IRepository {
	return &repo{pg: pg, red: red}
}

func (r *repo) Postgres() postgres.IRepository {
	return r.pg
}

func (r *repo) Redis() redis.IRepository {
	return r.red
}
