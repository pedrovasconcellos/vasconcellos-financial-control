package repository

import (
	"context"
	"time"
)

type ProcessedTransactionRepository interface {
	MarkProcessed(ctx context.Context, transactionID string, userID string, txnType string, processedAt time.Time) (bool, error)
	Remove(ctx context.Context, transactionID string) error
}
