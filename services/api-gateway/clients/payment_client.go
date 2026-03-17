package clients

import (
	"context"

	pb "ledgerflow/proto/paymentpb"

	"google.golang.org/grpc"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient() (*PaymentClient, error) {
	conn, err := grpc.Dial("payment-service:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewPaymentServiceClient(conn)

	return &PaymentClient{client: client, conn: conn}, nil
}

func (p *PaymentClient) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *PaymentClient) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest, opts ...grpc.CallOption) (*pb.CreatePaymentResponse, error) {
	return p.client.CreatePayment(ctx, req, opts...)
}
