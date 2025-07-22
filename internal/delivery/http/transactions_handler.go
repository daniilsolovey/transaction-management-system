package http

import (
	"log/slog"
	"net/http"

	kafka_broker "github.com/daniilsolovey/transaction-management-system/internal/delivery/kafka"
	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/daniilsolovey/transaction-management-system/internal/usecase"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type TransactionHandler struct {
	uc     *usecase.TransactionUseCase
	log    *slog.Logger
	writer *kafka_broker.Writer
}

func NewTransactionHandler(uc *usecase.TransactionUseCase, log *slog.Logger, writer *kafka_broker.Writer) *TransactionHandler {
	return &TransactionHandler{uc: uc, log: log, writer: writer}
}

func (h *TransactionHandler) RegisterRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/transactions", h.getTransactions)
	r.POST("/transactions", h.createTransaction)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

// getTransactions godoc
// @Summary Get transactions
// @Description Get transactions by user ID and optional type
// @Tags transactions
// @Produce json
// @Param user_id query string true "User ID"
// @Param type query string false "Transaction type (bet|win)"
// @Success 200 {array} domain.Transaction
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h *TransactionHandler) getTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	transactionType := c.Query("type") //  "bet" or "win"

	transactions, err := h.uc.Get(c.Request.Context(), userID, transactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if transactions == nil {
		transactions = []domain.Transaction{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, transactions)
}

// createTransaction godoc
// @Summary Create new transaction
// @Description Enqueue a new transaction via Kafka
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body domain.CreateTransactionMessage true "Transaction payload"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) createTransaction(c *gin.Context) {
	var message domain.CreateTransactionMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		h.log.Error("invalid request body", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Validate here before sending to Kafka
	if message.UserID == "" || message.Amount <= 0 ||
		(message.Type != "bet" && message.Type != "win") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction data"})
		return
	}

	// Send to Kafka
	err := h.writer.Publish(c.Request.Context(), message)
	if err != nil {
		h.log.Error("failed to publish to Kafka", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not enqueue message"})
		return
	}

	h.log.Info("transaction created", "user_id", message.UserID)
	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}
