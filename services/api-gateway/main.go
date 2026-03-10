package main

import (
	"ledgerflow/services/api-gateway/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	log.Println("Starting API Gateway on :8080")
	r.Run(":8080")
}

// POST /users
// POST /login
// POST /payments
//  GET /payments/:id
// POST /refunds
//  GET /transactions
