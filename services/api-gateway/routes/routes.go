package routes

import (
	"ledgerflow/services/api-gateway/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, paymentHandler *handlers.PaymentHandler) {
	r.POST("/payments", paymentHandler.CreatePayment)
}
