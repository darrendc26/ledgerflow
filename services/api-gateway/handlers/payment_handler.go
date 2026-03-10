package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	pb "ledgerflow/proto/paymentpb"
	"ledgerflow/services/api-gateway/clients"
)

func CreatePayment(c *gin.Context) {
	client := clients.NewPaymentClient()

	req := &pb.CreatePaymentRequest{
		UserId:   "user1",
		Amount:   "100",
		Currency: "USD",
	}

	resp, err := client.CreatePayment(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
