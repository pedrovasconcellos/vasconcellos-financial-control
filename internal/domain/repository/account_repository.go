package repository

import (
	"context"

	"github.com/vasconcellos/finance-control/internal/domain/entity"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id string, userID string) error
	GetByID(ctx context.Context, id string, userID string) (*entity.Account, error)
	List(ctx context.Context, userID string) ([]*entity.Account, error)
	AdjustBalance(ctx context.Context, id string, userID string, amount float64) error
}
