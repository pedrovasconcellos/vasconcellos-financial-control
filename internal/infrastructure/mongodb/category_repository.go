package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vasconcellos/finance-control/internal/domain/entity"
	domainErrors "github.com/vasconcellos/finance-control/internal/domain/errors"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type CategoryRepository struct {
	collection *mongo.Collection
}

var _ repository.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(client *Client) *CategoryRepository {
	return &CategoryRepository{collection: client.Collection("categories")}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	_, err := r.collection.InsertOne(ctx, category)
	return err
}

func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":     category.ID,
		"user_id": category.UserID,
	}, bson.M{"$set": bson.M{
		"name":        category.Name,
		"type":        category.Type,
		"description": category.Description,
		"parent_id":   category.ParentID,
		"updated_at":  time.Now().UTC(),
	}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string, userID string) error {
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

func (r *CategoryRepository) GetByID(ctx context.Context, id string, userID string) (*entity.Category, error) {
	var category entity.Category
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&category)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) List(ctx context.Context, userID string) ([]*entity.Category, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*entity.Category
	for cursor.Next(ctx) {
		var category entity.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}
