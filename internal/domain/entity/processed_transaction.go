package entity

import "time"

type ProcessedTransaction struct {
	ID            string    `bson:"_id"`
	TransactionID string    `bson:"transaction_id"`
	UserID        string    `bson:"user_id"`
	Type          string    `bson:"type"`
	ProcessedAt   time.Time `bson:"processed_at"`
}

