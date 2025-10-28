package repository

import (
	"context"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string, userID string) error
	GetByID(ctx context.Context, id string, userID string) (*entity.Category, error)
	List(ctx context.Context, userID string) ([]*entity.Category, error)
}
