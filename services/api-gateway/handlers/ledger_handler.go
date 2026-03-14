package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	pb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/api-gateway/clients"
)

type TransferRequest struct {
	SenderAccount   string `json:"sender_account" binding:"required"`
	ReceiverAccount string `json:"receiver_account" binding:"required"`
	Amount          int64  `json:"amount" binding:"required"`
	ReferenceId     string `json:"reference_id" binding:"required"`
}

type LedgerHandler struct {
	client *clients.LedgerClient
}

func NewLedgerHandler(client *clients.LedgerClient) *LedgerHandler {
	return &LedgerHandler{client: client}
}

func (h *LedgerHandler) Transfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pbReq := &pb.TransferRequest{
		SenderAccount:   req.SenderAccount,
		ReceiverAccount: req.ReceiverAccount,
		Amount:          req.Amount,
		ReferenceId:     req.ReferenceId,
	}

	resp, err := h.client.Transfer(context.Background(), pbReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         resp.Status,
		"transaction_id": resp.TransactionId,
	})
}
