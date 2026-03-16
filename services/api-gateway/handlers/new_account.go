package handlers

import (
	"net/http"

	"ledgerflow/infra/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewAccountRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type AccountHandler struct {
	db      *pgxpool.Pool
	metrics *prometheus.Prometheus
}

func NewAccountHandler(db *pgxpool.Pool, metrics *prometheus.Prometheus) *AccountHandler {
	return &AccountHandler{db: db, metrics: metrics}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {

	var req NewAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var accountID string

	query := `
	INSERT INTO accounts (user_id, balance, currency)
	VALUES ($1, 0, $2)
	RETURNING id
	`

	err := h.db.QueryRow(c.Request.Context(), query, req.UserID, req.Currency).Scan(&accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.metrics.AccountsCreated.Inc()

	c.JSON(http.StatusOK, gin.H{
		"account_id": accountID,
		"user_id":    req.UserID,
		"status":     "created",
	})
}
