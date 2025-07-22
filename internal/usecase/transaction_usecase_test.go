package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/postgres"
	"github.com/daniilsolovey/transaction-management-system/internal/repository/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ postgres.IRepository = (*mockPostgresRepo)(nil)

type mockPostgresRepo struct {
	mock.Mock
}

type mockRepo struct {
	pg    postgres.IRepository
	redis redis.IRepository
}

func (m *mockPostgresRepo) SaveTransaction(ctx context.Context,
	tx domain.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *mockPostgresRepo) GetFilteredTransactions(ctx context.Context,
	userID, txType string) ([]domain.Transaction, error) {

	args := m.Called(ctx, userID, txType)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *mockPostgresRepo) Close() {
}

func (m *mockRepo) Postgres() postgres.IRepository {
	return m.pg
}

func (m *mockRepo) Redis() redis.IRepository {
	return m.redis
}

func TestCreateTransaction_Valid(t *testing.T) {
	mockPg := new(mockPostgresRepo)
	mockRepo := &mockRepo{pg: mockPg}
	logger := slog.Default()
	uc := NewTransactionUseCase(mockRepo, logger)

	txMsg := domain.CreateTransactionMessage{
		UserID: "123e4567-e89b-12d3-a456-426614174000",
		Type:   "bet",
		Amount: 50.0,
	}

	mockPg.On("SaveTransaction",
		mock.Anything, mock.AnythingOfType("domain.Transaction")).Return(nil)

	err := uc.Create(context.Background(), txMsg)
	assert.NoError(t, err)
	mockPg.AssertExpectations(t)
}

func TestCreateTransaction_InvalidType(t *testing.T) {
	mockPg := new(mockPostgresRepo)
	mockRepo := &mockRepo{pg: mockPg}
	logger := slog.Default()
	uc := NewTransactionUseCase(mockRepo, logger)

	txMsg := domain.CreateTransactionMessage{
		UserID: "abc",
		Type:   "invalid",
		Amount: 100,
	}

	err := uc.Create(context.Background(), txMsg)
	assert.EqualError(t, err, "transaction type must be 'bet' or 'win'")
	mockPg.AssertNotCalled(t, "SaveTransaction")
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	mockPg := new(mockPostgresRepo)
	mockRepo := &mockRepo{pg: mockPg}
	logger := slog.Default()
	uc := NewTransactionUseCase(mockRepo, logger)

	txMsg := domain.CreateTransactionMessage{
		UserID: "abc",
		Type:   "bet",
		Amount: 0,
	}

	err := uc.Create(context.Background(), txMsg)
	assert.EqualError(t, err, "amount must be greater than 0")
	mockPg.AssertNotCalled(t, "SaveTransaction")
}

func TestGetTransactions(t *testing.T) {
	mockPg := new(mockPostgresRepo)
	mockRepo := &mockRepo{pg: mockPg}
	logger := slog.Default()
	uc := NewTransactionUseCase(mockRepo, logger)

	expected := []domain.Transaction{
		{UserID: "abc", Type: "bet", Amount: 10.5, Timestamp: time.Now()},
	}

	mockPg.On("GetFilteredTransactions", mock.Anything,
		"abc", "bet").Return(expected, nil)

	result, err := uc.Get(context.Background(), "abc", "bet")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockPg.AssertExpectations(t)
}

func TestGetTransactions_Error(t *testing.T) {
	mockPg := new(mockPostgresRepo)
	mockRepo := &mockRepo{pg: mockPg}
	logger := slog.Default()
	uc := NewTransactionUseCase(mockRepo, logger)

	mockPg.On("GetFilteredTransactions", mock.Anything, "abc", "bet").
		Return(([]domain.Transaction)(nil), errors.New("db error"))

	result, err := uc.Get(context.Background(), "abc", "bet")
	assert.Error(t, err)
	assert.Nil(t, result)
}
