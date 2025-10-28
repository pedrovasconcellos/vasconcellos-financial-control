package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
	domainErrors "github.com/vasconcellos/financial-control/src/internal/domain/errors"
	"github.com/vasconcellos/financial-control/src/internal/domain/repository"
)

type AccountRepository struct {
	collection *mongo.Collection
}

var _ repository.AccountRepository = (*AccountRepository)(nil)

func NewAccountRepository(client *Client) (*AccountRepository, error) {
	col := client.Collection("accounts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}
	if _, err := col.Indexes().CreateMany(ctx, indexModels); err != nil {
		return nil, err
	}
	return &AccountRepository{collection: col}, nil
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	_, err := r.collection.InsertOne(ctx, account)
	if mongo.IsDuplicateKeyError(err) {
		return domainErrors.ErrConflict
	}
	return err
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     account.ID,
		"user_id": account.UserID,
	}, bson.M{"$set": bson.M{
		"name":        account.Name,
		"type":        account.Type,
		"currency":    account.Currency,
		"description": account.Description,
		"updated_at":  account.UpdatedAt,
	}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id string, userID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id string, userID string) (*entity.Account, error) {
	var account entity.Account
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&account)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) List(ctx context.Context, userID string, limit int64, offset int64) ([]*entity.Account, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if offset > 0 {
		opts.SetSkip(offset)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*entity.Account
	for cursor.Next(ctx) {
		var account entity.Account
		if err := cursor.Decode(&account); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (r *AccountRepository) AdjustBalance(ctx context.Context, id string, userID string, amount float64) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}, bson.M{
		"$inc": bson.M{"balance": amount},
		"$set": bson.M{"updated_at": time.Now().UTC()},
	})

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}
