package domain

import "time"

// Transaction represents a stored transaction record.
type Transaction struct {
	UserID    string
	Type      string // "bet" or "win"
	Amount    float64
	Timestamp time.Time
}

// CreateTransactionMessage is the input data for creating a new transaction.
// This may come from a Kafka message or an HTTP API payload.
type CreateTransactionMessage struct {
	UserID    string    `json:"user_id"`
	Type      string    `json:"transaction_type"` // "bet" or "win"
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"` // Optional: can be filled in if not present
}
