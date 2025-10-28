package repository

import (
	"context"
	"time"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
)

type BudgetRepository interface {
	Create(ctx context.Context, budget *entity.Budget) error
	Update(ctx context.Context, budget *entity.Budget) error
	GetByID(ctx context.Context, id string, userID string) (*entity.Budget, error)
	List(ctx context.Context, userID string, limit int64, offset int64) ([]*entity.Budget, error)
	UpdateSpent(ctx context.Context, id string, userID string, spent float64) error
	FindActiveByCategory(ctx context.Context, userID string, categoryID string, timestamp time.Time) ([]*entity.Budget, error)
}
