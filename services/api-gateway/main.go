package main

import (
	"ledgerflow/services/api-gateway/clients"
	"ledgerflow/services/api-gateway/handlers"
	"ledgerflow/services/api-gateway/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	paymentClient, err := clients.NewPaymentClient()
	if err != nil {
		log.Fatalf("failed to initialize payment client: %v", err)
	}
	defer paymentClient.Close()

	paymentHandler := handlers.NewPaymentHandler(paymentClient)

	r := gin.Default()

	routes.RegisterRoutes(r, paymentHandler)

	log.Println("Starting API Gateway on :8080")
	r.Run(":8080")
}

// POST /users
// POST /login
// POST /payments
//  GET /payments/:id
// POST /refunds
//  GET /transactions
