package routes

import (
	"ledgerflow/services/api-gateway/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/payments", handlers.CreatePayment)
}
