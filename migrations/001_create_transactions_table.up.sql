CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('bet', 'win')),
    amount NUMERIC(12, 2) NOT NULL CHECK (amount >= 0),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
