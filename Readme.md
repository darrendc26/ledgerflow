# LedgerFlow
LedgerFlow is a distributed payment processing system designed to demonstrate production-grade backend architecture. The project implements a microservice-based payment infrastructure with event-driven processing, financial ledger integrity, observability, and reliability mechanisms commonly used in real-world fintech systems.

The goal of this project is to showcase backend engineering skills relevant to modern infrastructure teams working with high-scale distributed systems.

## Overview
LedgerFlow simulates a simplified payment processing platform where users can create accounts, perform transfers, and have transactions processed asynchronously through an event-driven architecture.

The system includes:
- REST API gateway
- gRPC microservices
- event-driven architecture using Kafka
- asynchronous workers
- financial double-entry ledger
- retry and dead-letter queue handling
- distributed tracing
- metrics and monitoring
- observability dashboards

This architecture mirrors patterns used in production systems at large technology companies.

## Architecture
            Client
              │
              ▼
    API Gateway (REST - Gin)
              │
              ▼
     Payment Service (gRPC)
              │
              ▼
         Kafka Event Bus
              │
              ▼
         Payment Worker
              │
              ▼
      Ledger Service (gRPC)
              │
              ▼
           PostgreSQL

## Observability layer
Application Metrics → Prometheus  
Dashboards → Grafana  
Tracing → OpenTelemetry  
Trace Visualization → Jaeger  

## Core Components

### API Gateway
Handles incoming HTTP requests and routes them to backend services.

Responsibilities:
- account creation
- initiating payments
- request validation
- communication with payment service via gRPC

Technology stack:
- Go
- Gin web framework

### Payment Service
Handles payment creation and orchestrates the payment workflow.

Responsibilities:
- create payment records
- publish payment events to Kafka
- maintain payment status lifecycle

Payment status states:
- created
- processing
- completed
- failed

### Ledger Service
Implements a financial ledger using double-entry accounting.

Responsibilities:
- debit sender account
- credit receiver account
- ensure transactional integrity
- maintain account balances

Every transaction creates two ledger entries ensuring sum(debits) = sum(credits)

### Payment Worker
Processes payment events asynchronously from Kafka.

Responsibilities:
- consume payment events
- perform ledger transfers
- retry failed operations
- publish failed events to dead-letter queue
- update payment status
- emit metrics and traces

Retry policy:
- three retry attempts
- exponential delay
- DLQ fallback on persistent failure

### Event Bus
Kafka is used for decoupled event-driven processing.

Topics used:
payments  
payments_dlq  

Dead-letter queue allows recovery and investigation of failed transactions.

## Data Model

### Accounts
id  
user_id  
balance  
currency  

### Payments
id  
sender_account  
receiver_account  
amount  
currency  
status  
created_at  

### Ledger Entries
id  
transaction_id  
account_id  
amount  
type  

Ledger entry types:
- debit
- credit

## Observability
LedgerFlow implements full observability across metrics, tracing, and logs.

### Metrics
Prometheus metrics include:
payments_processed_total
payments_failed_total
payment_processing_seconds
accounts_created_total

Histogram metrics track latency distribution.

Example metric:
payment_processing_seconds_bucket

### Monitoring
Grafana dashboards visualize:
- payment throughput
- failure rate
- latency percentiles
- total transactions processed

Example queries:
rate(payments_processed_total[1m])  
rate(payments_failed_total[5m])  
histogram_quantile(0.95, rate(payment_processing_seconds_bucket[5m]))  
accounts_created_total  

### Distributed Tracing
OpenTelemetry instrumentation tracks requests across services.  

Trace example:  
process-payment  
├── kafka-consume  
├── ledger-transfer  
└── status-update  

Jaeger provides a visualization interface for tracing.  

## Reliability Features
The system includes several reliability mechanisms.  

### Idempotency
Duplicate events are ignored using payment status validation.

### Retry Handling
Transient failures are retried automatically.

### Dead Letter Queue
Persistent failures are redirected to a dedicated Kafka topic.

### Financial Integrity
Double-entry ledger guarantees accounting correctness.

## Infrastructure Stack
Language:  
Go  

Infrastructure:  
Kafka  
PostgreSQL  
Docker  

Observability:  
Prometheus  
Grafana  
OpenTelemetry  
Jaeger  
  
Communication:  
REST  
gRPC  

## Running the System
### Prerequisites:  
Docker  
Docker Compose  
Go  

### Setup
```bash
# Clone repository
git clone https://github.com/darrendc26/ledgerflow.git
cd ledgerflow

# Build and Run 
docker compose up --build  
```

### Testing endpoints with curl 
```bash
# Create accounts
curl -X POST http://localhost:8080/accounts \
-H "Content-Type: application/json" \
-d '{
  "user_id": "user_1",
  "currency": "USD"
}'

curl -X POST http://localhost:8080/accounts \
-H "Content-Type: application/json" \
-d '{
  "user_id": "user_2",
  "currency": "USD"
}'

# Deposit
curl -X POST http://localhost:8080/deposits \
-H "Content-Type: application/json" \
-d '{
  "deposit_account": "10000000",
  "amount": 1000,
  "currency": "USD"
}'

# Transfer funds
curl -X POST http://localhost:8080/payments \
-H "Content-Type: application/json" \
-d '{
  "sender_account": "10000000",
  "receiver_account": "10000001",
  "amount": 100,
  "currency": "USD"
}'

# Testing DLQ (dead letter queue) for failed payments
curl -X POST http://localhost:8080/payments \
-H "Content-Type: application/json" \
-d '{
  "sender_account": "10000000",
  "receiver_account": "10000001",
  "amount": 100000,
  "currency": "USD"
}'

# Load test
for i in {1..20}; do
curl -X POST http://localhost:8080/payments \
-H "Content-Type: application/json" \
-d '{
  "sender_account": "10000000",
  "receiver_account": "10000001",
  "amount": 1,
  "currency": "USD"
}'
done
```

## Services
| Service       | URL                     |
|--------------|--------------------------|
| API Gateway  | http://localhost:8080    |
| Prometheus   | http://localhost:9090    |
| Grafana      | http://localhost:3000    |
| Jaeger       | http://localhost:16686   |

## What This Project Demonstrates
- event-driven architecture
- distributed microservices
- asynchronous job processing
- financial ledger design
- reliability and failure handling
- distributed tracing
- metrics-driven monitoring
- production-style observability

## Future Improvements
- ledger reconciliation service
- idempotency keys for exactly-once semantics
- authentication and authorization
- rate limiting
- automated alerting
- horizontal worker scaling
- multi-currency support

## Author
Darren Da Costa