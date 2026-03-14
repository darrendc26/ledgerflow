package server

import (
	"log"
	"net"

	pb "ledgerflow/proto/ledgerpb"
	"ledgerflow/services/ledger/handler"

	"google.golang.org/grpc"
)

func StartGrpcServer(handler *handler.LedgerHandler) {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterLedgerServiceServer(grpcServer, handler)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
