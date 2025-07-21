package repository

import (
	"github.com/daniilsolovey/transaction-management-system/internal/repository/postgres"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/redis"
)

type IRepository interface {
	Postgres() postgres.IRepository
	Redis() redis.IRepository
}
