package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

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

func (uc *TransactionUseCase) Create(ctx context.Context,
	message domain.CreateTransactionMessage) error {
	uc.log.Info("creating transaction",
		"user_id", message.UserID,
		"type", message.Type,
		"amount", message.Amount,
		"timestamp", message.Timestamp,
	)

	if message.UserID == "" {
		err := errors.New("user_id is required")
		uc.log.Error("validation failed", "err", err)
		return err
	}

	if message.Type != "bet" && message.Type != "win" {
		err := errors.New("transaction type must be 'bet' or 'win'")
		uc.log.Error("validation failed", "err", err)
		return err
	}

	if message.Amount <= 0 {
		err := errors.New("amount must be greater than 0")
		uc.log.Error("validation failed", "err", err)
		return err
	}

	timestamp := message.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}

	transaction := domain.Transaction{
		UserID:    message.UserID,
		Type:      message.Type,
		Amount:    message.Amount,
		Timestamp: timestamp,
	}

	if err := uc.repo.Postgres().SaveTransaction(ctx, transaction); err != nil {
		uc.log.Error("failed to save transaction", "err", err, "tx", transaction)
		return err
	}

	uc.log.Info("transaction saved successfully", "user_id", transaction.UserID)
	return nil
}

func (uc *TransactionUseCase) Get(ctx context.Context,
	userID, transactionType string) ([]domain.Transaction, error) {
	uc.log.Info("retrieving transactions", "user_id", userID, "type", transactionType)
	transactions, err := uc.repo.Postgres().GetFilteredTransactions(ctx, userID, transactionType)
	if err != nil {
		uc.log.Error("failed to retrieve transactions", "err", err)
		return nil, err
	}

	uc.log.Info("transactions retrieved", "count", len(transactions))
	return transactions, nil
}
