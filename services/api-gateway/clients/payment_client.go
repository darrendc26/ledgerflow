package clients

import (
	"log"

	pb "ledgerflow/proto/paymentpb"

	"google.golang.org/grpc"
)

func NewPaymentClient() pb.PaymentServiceClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return pb.NewPaymentServiceClient(conn)
}
