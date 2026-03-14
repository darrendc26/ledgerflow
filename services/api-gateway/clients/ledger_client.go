package clients

import (
	"context"

	pb "ledgerflow/proto/ledgerpb"

	"google.golang.org/grpc"
)

type LedgerClient struct {
	client pb.LedgerServiceClient
	conn   *grpc.ClientConn
}

func NewLedgerClient() (*LedgerClient, error) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewLedgerServiceClient(conn)

	return &LedgerClient{client: client, conn: conn}, nil
}

func (l *LedgerClient) Close() error {
	if l.conn != nil {
		return l.conn.Close()
	}
	return nil
}

func (l *LedgerClient) Transfer(ctx context.Context, req *pb.TransferRequest, opts ...grpc.CallOption) (*pb.TransferResponse, error) {
	return l.client.Transfer(ctx, req, opts...)
}
