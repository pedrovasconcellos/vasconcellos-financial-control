package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
	domainErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
)

type BudgetRepository struct {
	collection *mongo.Collection
}

var _ repository.BudgetRepository = (*BudgetRepository)(nil)

func NewBudgetRepository(client *Client) *BudgetRepository {
	return &BudgetRepository{collection: client.Collection("budgets")}
}

func (r *BudgetRepository) Create(ctx context.Context, budget *entity.Budget) error {
	_, err := r.collection.InsertOne(ctx, budget)
	return err
}

func (r *BudgetRepository) Update(ctx context.Context, budget *entity.Budget) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     budget.ID,
		"user_id": budget.UserID,
	}, bson.M{"$set": budget})
	return err
}

func (r *BudgetRepository) GetByID(ctx context.Context, id string, userID string) (*entity.Budget, error) {
	var budget entity.Budget
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&budget)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *BudgetRepository) List(ctx context.Context, userID string) ([]*entity.Budget, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []*entity.Budget
	for cursor.Next(ctx) {
		var budget entity.Budget
		if err := cursor.Decode(&budget); err != nil {
			return nil, err
		}
		budgets = append(budgets, &budget)
	}
	return budgets, nil
}

func (r *BudgetRepository) UpdateSpent(ctx context.Context, id string, userID string, spent float64) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}, bson.M{"$set": bson.M{"spent": spent}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *BudgetRepository) FindActiveByCategory(ctx context.Context, userID string, categoryID string, timestamp time.Time) ([]*entity.Budget, error) {
	filter := bson.M{
		"user_id":      userID,
		"category_id":  categoryID,
		"period_start": bson.M{"$lte": timestamp},
		"period_end":   bson.M{"$gte": timestamp},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []*entity.Budget
	for cursor.Next(ctx) {
		var budget entity.Budget
		if err := cursor.Decode(&budget); err != nil {
			return nil, err
		}
		budgets = append(budgets, &budget)
	}

	return budgets, nil
}
