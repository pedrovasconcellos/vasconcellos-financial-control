package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
	"github.com/vasconcellos/financial-control/src/internal/domain/repository"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

var _ repository.TransactionRepository = (*TransactionRepository)(nil)

func NewTransactionRepository(client *Client) (*TransactionRepository, error) {
	col := client.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "occurred_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "category_id", Value: 1},
				{Key: "occurred_at", Value: -1},
			},
		},
	}
	if _, err := col.Indexes().CreateMany(ctx, indexModels); err != nil {
		return nil, err
	}

	return &TransactionRepository{collection: col}, nil
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     transaction.ID,
		"user_id": transaction.UserID,
	}, bson.M{"$set": bson.M{
		"category_id":    transaction.CategoryID,
		"description":    transaction.Description,
		"notes":          transaction.Notes,
		"tags":           transaction.Tags,
		"status":         transaction.Status,
		"receipt_object": transaction.ReceiptObject,
		"metadata":       transaction.Metadata,
		"updated_at":     time.Now().UTC(),
	}})
	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string, userID string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&transaction)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) List(ctx context.Context, userID string, from time.Time, to time.Time, limit int64, offset int64) ([]*entity.Transaction, error) {
	filter := bson.M{
		"user_id": userID,
		"occurred_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "occurred_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if offset > 0 {
		opts.SetSkip(offset)
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*entity.Transaction
	for cursor.Next(ctx) {
		var transaction entity.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) ListByCategory(ctx context.Context, userID string, categoryID string, from time.Time, to time.Time) ([]*entity.Transaction, error) {
	filter := bson.M{
		"user_id":     userID,
		"category_id": categoryID,
		"occurred_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*entity.Transaction
	for cursor.Next(ctx) {
		var transaction entity.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}
