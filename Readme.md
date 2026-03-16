Testing endpoints with curl:
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