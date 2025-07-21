package postgres

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	sql  squirrel.StatementBuilderType
	log  *slog.Logger
}

type IRepository interface {
	SaveTransaction(ctx context.Context, transaction domain.Transaction) error
	GetFilteredTransactions(ctx context.Context, userID string, transactionType string) ([]domain.Transaction, error)
	Close()
}

func New(pool *pgxpool.Pool, logger *slog.Logger) *Repository {
	return &Repository{
		pool: pool,
		sql:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:  logger,
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *Repository) Close() {
	r.pool.Close()
}
