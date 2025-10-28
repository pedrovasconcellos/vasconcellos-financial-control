package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
	"github.com/vasconcellos/financial-control/src/internal/domain/repository"
)

type ProcessedTransactionRepository struct {
	collection *mongo.Collection
}

var _ repository.ProcessedTransactionRepository = (*ProcessedTransactionRepository)(nil)

func NewProcessedTransactionRepository(client *Client) *ProcessedTransactionRepository {
	col := client.Collection("processed_transactions")
	return &ProcessedTransactionRepository{collection: col}
}

func (r *ProcessedTransactionRepository) MarkProcessed(ctx context.Context, transactionID string, userID string, txnType string, processedAt time.Time) (bool, error) {
	doc := entity.ProcessedTransaction{
		ID:            transactionID,
		TransactionID: transactionID,
		UserID:        userID,
		Type:          txnType,
		ProcessedAt:   processedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ProcessedTransactionRepository) Remove(ctx context.Context, transactionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": transactionID})
	return err
}
