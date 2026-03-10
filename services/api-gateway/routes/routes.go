package routes

import (
	"github.com/gin-gonic/gin"

	"ledgerflow/services/api-gateway/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/payments", handlers.CreatePayment)
}
