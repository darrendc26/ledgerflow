FROM golang:1.25-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ledger-service ./services/ledger
RUN go build -o payment-service ./services/payment-service
RUN go build -o api-gateway ./services/api-gateway
RUN go build -o payment-worker ./services/payment-worker