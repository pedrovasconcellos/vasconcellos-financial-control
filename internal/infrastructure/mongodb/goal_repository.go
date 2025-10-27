package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
	domainErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
)

type GoalRepository struct {
	collection *mongo.Collection
}

var _ repository.GoalRepository = (*GoalRepository)(nil)

func NewGoalRepository(client *Client) *GoalRepository {
	return &GoalRepository{collection: client.Collection("goals")}
}

func (r *GoalRepository) Create(ctx context.Context, goal *entity.Goal) error {
	_, err := r.collection.InsertOne(ctx, goal)
	return err
}

func (r *GoalRepository) Update(ctx context.Context, goal *entity.Goal) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     goal.ID,
		"user_id": goal.UserID,
	}, bson.M{"$set": goal})
	return err
}

func (r *GoalRepository) GetByID(ctx context.Context, id string, userID string) (*entity.Goal, error) {
	var goal entity.Goal
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&goal)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *GoalRepository) List(ctx context.Context, userID string) ([]*entity.Goal, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var goals []*entity.Goal
	for cursor.Next(ctx) {
		var goal entity.Goal
		if err := cursor.Decode(&goal); err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}

func (r *GoalRepository) UpdateProgress(ctx context.Context, id string, userID string, amount float64) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}, bson.M{"$inc": bson.M{"current_amount": amount}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}
