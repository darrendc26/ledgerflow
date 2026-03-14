package routes

import (
	"ledgerflow/services/api-gateway/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	paymentHandler *handlers.PaymentHandler,
	NewAccountHandler *handlers.AccountHandler,
	depositHandler *handlers.DepositHandler,

) {
	r.POST("/payments", paymentHandler.CreatePayment)
	r.POST("/accounts", NewAccountHandler.CreateAccount)
	r.POST("/deposits", depositHandler.CreateDeposit)
}
