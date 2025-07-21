package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/daniilsolovey/transaction-management-system/internal/domain"
)

func (r *Repository) SaveTransaction(ctx context.Context, transaction domain.Transaction) error {
	query, args, err := r.sql.
		Insert("transactions").
		Columns("user_id", "transaction_type", "amount", "timestamp").
		Values(transaction.UserID, transaction.Type, transaction.Amount, transaction.Timestamp).
		ToSql()
	if err != nil {
		r.log.Error("failed to build insert", "err", err)
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		r.log.Error("failed to insert transaction", "err", err)
	}
	return err
}

func (r *Repository) GetFilteredTransactions(ctx context.Context,
	userID string, transactionType string) ([]domain.Transaction, error) {
	builder := r.sql.
		Select("user_id", "transaction_type", "amount", "timestamp").
		From("transactions").
		OrderBy("timestamp DESC")

	// Filter by user ID if provided
	if userID != "" {
		builder = builder.Where(squirrel.Eq{"user_id": userID})
	}

	// Normalize and apply transaction type filter
	switch transactionType {
	case "bet", "win":
		builder = builder.Where(squirrel.Eq{"transaction_type": transactionType})
	case "", "all":
		// No filtering
	default:
		r.log.Warn("unsupported transaction type filter", "type", transactionType)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build select query", "err", err)
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		r.log.Error("failed to query transactions", "err", err)
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		err := rows.Scan(&tx.UserID, &tx.Type, &tx.Amount, &tx.Timestamp)
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}

	return result, nil
}
