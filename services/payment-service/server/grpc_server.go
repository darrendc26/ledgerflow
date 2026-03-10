package server

import (
	"log"
	"net"

	pb "ledgerflow/proto/paymentpb"
	"ledgerflow/services/payment-service/handler"

	"google.golang.org/grpc"
)

func StartGrpcServer(handler *handler.PaymentHandler) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, handler)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
