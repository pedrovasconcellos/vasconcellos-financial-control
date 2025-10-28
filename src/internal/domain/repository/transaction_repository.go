package repository

import (
	"context"
	"time"

	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	Update(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id string, userID string) (*entity.Transaction, error)
	List(ctx context.Context, userID string, from time.Time, to time.Time, limit int64, offset int64) ([]*entity.Transaction, error)
	ListByCategory(ctx context.Context, userID string, categoryID string, from time.Time, to time.Time) ([]*entity.Transaction, error)
}
