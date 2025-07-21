curl -X POST http://localhost:3000/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "transaction_type": "bet",
    "amount": 101.50,
    "timestamp": "2025-07-21T14:30:00Z"
}'
