package entity

import "time"

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID            string            `bson:"_id"`
	UserID        string            `bson:"user_id"`
	AccountID     string            `bson:"account_id"`
	CategoryID    string            `bson:"category_id"`
	Amount        float64           `bson:"amount"`
	Currency      Currency          `bson:"currency"`
	Description   string            `bson:"description"`
	OccurredAt    time.Time         `bson:"occurred_at"`
	Status        TransactionStatus `bson:"status"`
	Notes         string            `bson:"notes"`
	ReceiptObject *string           `bson:"receipt_object,omitempty"`
	Tags          []string          `bson:"tags"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     time.Time         `bson:"updated_at"`
	ExternalRef   string            `bson:"external_ref"`
	Metadata      map[string]string `bson:"metadata"`
}
