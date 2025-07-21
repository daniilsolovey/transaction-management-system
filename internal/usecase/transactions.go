package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/daniilsolovey/transaction-management-system/internal/repository"
)

type TransactionRepo interface {
	Save(ctx context.Context, tx *domain.Transaction) error
	GetFiltered(ctx context.Context, userID, txType string) ([]domain.Transaction, error)
}

type TransactionUseCase struct {
	repo repository.IRepository
	log  *slog.Logger
}

func NewTransactionUseCase(repo repository.IRepository, log *slog.Logger) *TransactionUseCase {
	return &TransactionUseCase{repo: repo, log: log}
}

func (uc *TransactionUseCase) Create(ctx context.Context, message domain.CreateTransactionMessage) error {
	uc.log.Info("creating transaction", "user_id", message.UserID, "type", message.Type, "amount", message.Amount)

	if message.UserID == "" || message.Amount <= 0 {
		err := errors.New("invalid transaction data")
		uc.log.Error("validation failed", "err", err, "dto", message)
		return err
	}

	transaction := domain.Transaction{
		UserID:    message.UserID,
		Type:      message.Type,
		Amount:    message.Amount,
		Timestamp: message.Timestamp, // make sure dto includes Timestamp or set it here
	}

	if err := uc.repo.Postgres().SaveTransaction(ctx, transaction); err != nil {
		uc.log.Error("failed to save transaction", "err", err, "tx", transaction)
		return err
	}

	uc.log.Info("transaction saved successfully", "user_id", transaction.UserID)
	return nil
}

func (uc *TransactionUseCase) Get(ctx context.Context, userID, txType string) ([]domain.Transaction, error) {
	uc.log.Info("retrieving transactions", "user_id", userID, "type", txType)
	transactions, err := uc.repo.Postgres().GetFilteredTransactions(ctx, userID, txType)
	if err != nil {
		uc.log.Error("failed to retrieve transactions", "err", err)
		return nil, err
	}

	uc.log.Info("transactions retrieved", "count", len(transactions))
	return transactions, nil
}
