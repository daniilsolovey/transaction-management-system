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

func (r *Repository) GetFilteredTransactions(ctx context.Context, userID string, txType string) ([]domain.Transaction, error) {
	builder := r.sql.
		Select("user_id", "transaction_type", "amount", "timestamp").
		From("transactions").
		OrderBy("timestamp DESC")

	if userID != "" {
		builder = builder.Where(squirrel.Eq{"user_id": userID})
	}

	if txType != "" && txType != "all" {
		builder = builder.Where(squirrel.Eq{"transaction_type": txType})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build select", "err", err)
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
		if err := rows.Scan(&tx.UserID, &tx.Type, &tx.Amount, &tx.Timestamp); err != nil {
			return nil, err
		}
		result = append(result, tx)
	}

	return result, nil
}
