package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type PaymentEvent struct {
	PaymentID       string
	SenderAccount   string
	ReceiverAccount string
	Amount          int64
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP("localhost:9092"),
			Topic: "payments",
		},
	}
}

func (p *Producer) Publish(paymentEvent *PaymentEvent) error {
	data, err := json.Marshal(paymentEvent)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: "payments",
			Value: data,
		},
	)
}

func (p *Producer) PublishDLQ(paymentEvent *PaymentEvent) error {
	data, _ := json.Marshal(paymentEvent)

	return p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: "payments_dlq",
			Value: data,
		},
	)
}
