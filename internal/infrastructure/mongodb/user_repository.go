package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vasconcellos/finance-control/internal/domain/entity"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type UserRepository struct {
	collection *mongo.Collection
}

var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository(client *Client) (*UserRepository, error) {
	col := client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, index); err != nil {
		return nil, err
	}
	return &UserRepository{collection: col}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := r.collection.UpdateByID(ctx, user.ID, bson.M{"$set": user})
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
