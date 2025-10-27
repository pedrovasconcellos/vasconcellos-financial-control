package repository

import (
	"context"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *entity.Goal) error
	Update(ctx context.Context, goal *entity.Goal) error
	GetByID(ctx context.Context, id string, userID string) (*entity.Goal, error)
	List(ctx context.Context, userID string) ([]*entity.Goal, error)
	UpdateProgress(ctx context.Context, id string, userID string, amount float64) error
}
