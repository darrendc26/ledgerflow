package main

import (
	"context"
	"log"
	"net"

	pb "ledgerflow/proto/paymentpb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *server) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {

	log.Println("CreatePayment caled")

	return &pb.CreatePaymentResponse{
		PaymentId: "1234",
		Status:    "success",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterPaymentServiceServer(grpcServer, &server{})

	log.Println("Payment Service running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
