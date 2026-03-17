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
	paymentWriter *kafka.Writer
	dlqWriter     *kafka.Writer
}

func NewProducer() *Producer {

	return &Producer{
		paymentWriter: &kafka.Writer{
			Addr:  kafka.TCP("kafka:9092"),
			Topic: "payments",
		},

		dlqWriter: &kafka.Writer{
			Addr:  kafka.TCP("kafka:9092"),
			Topic: "payments_dlq",
		},
	}
}

func (p *Producer) Publish(event *PaymentEvent) error {

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.paymentWriter.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: data,
		},
	)
}

func (p *Producer) PublishDLQ(event *PaymentEvent) error {

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.dlqWriter.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: data,
		},
	)
}
