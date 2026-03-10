package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	pb "ledgerflow/proto/paymentpb"
	"ledgerflow/services/api-gateway/clients"
)

type CreatePaymentRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := clients.NewPaymentClient()

	pbReq := &pb.CreatePaymentRequest{
		UserId:   req.UserID,
		Amount:   req.Amount,
		Currency: req.Currency,
	}

	resp, err := client.CreatePayment(context.Background(), pbReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
