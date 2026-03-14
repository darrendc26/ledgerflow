package handlers

import (
	"net/http"

	"ledgerflow/services/api-gateway/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateDepositRequest struct {
	DepositAccount string `json:"deposit_account" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
	Currency       string `json:"currency" binding:"required"`
}

type DepositHandler struct {
	db *pgxpool.Pool
}

func NewDepositHandler(db *pgxpool.Pool) *DepositHandler {
	return &DepositHandler{db: db}
}

func (h *DepositHandler) CreateDeposit(c *gin.Context) {

	var req CreateDepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// txID := uuid.New().String()

	// _, err := h.db.Query(c.Request.Context(),
	// 	`INSERT INTO transactions (id) VALUES ($1)`,
	// 	txID,
	// )

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	deposit_repo := repo.NewDepositRepository(h.db)
	txID, err := deposit_repo.CreateDeposit(c.Request.Context(), req.DepositAccount, req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": txID,
		"account_id":     req.DepositAccount,
		"amount":         req.Amount,
		"status":         "deposit",
	})
}
