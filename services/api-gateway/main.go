package main

import (
	"log"

	"github.com/gin-gonic/gin"

	// "ledgerflow/infra/prometheus"
	"ledgerflow/pkg/db"
	"ledgerflow/services/api-gateway/clients"
	"ledgerflow/services/api-gateway/handlers"
	"ledgerflow/services/api-gateway/routes"
)

func main() {

	pool := db.NewPostgresPool()
	accountHandler := handlers.NewAccountHandler(pool)
	depositHandler := handlers.NewDepositHandler(pool)
	// payment client
	paymentClient, err := clients.NewPaymentClient()
	if err != nil {
		log.Fatalf("failed to initialize payment client: %v", err)
	}
	defer paymentClient.Close()

	paymentHandler := handlers.NewPaymentHandler(paymentClient)

	r := gin.Default()

	routes.RegisterRoutes(r, paymentHandler, accountHandler, depositHandler)

	log.Println("Starting API Gateway on :8080")
	r.Run(":8080")
}

// POST /users
// POST /login
// POST /payments
//  GET /payments/:id
// POST /refunds
//  GET /transactions
