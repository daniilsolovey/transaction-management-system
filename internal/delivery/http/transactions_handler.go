package http

import (
	"log/slog"
	"net/http"

	"github.com/daniilsolovey/transaction-management-system/internal/domain"
	"github.com/daniilsolovey/transaction-management-system/internal/usecase"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type TransactionHandler struct {
	uc  *usecase.TransactionUseCase
	log *slog.Logger
}

func NewTransactionHandler(uc *usecase.TransactionUseCase, log *slog.Logger) *TransactionHandler {
	return &TransactionHandler{uc: uc, log: log}
}

func (h *TransactionHandler) RegisterRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/transactions", h.getTransactions)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

// getTransactions handles requests to query transactions.
func (h *TransactionHandler) getTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	txType := c.Query("type") // e.g., "bet" or "win"

	transactions, err := h.uc.Get(c.Request.Context(), userID, txType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if transactions == nil {
		transactions = []domain.Transaction{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, transactions)
}
