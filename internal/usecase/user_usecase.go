package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) EnsureUser(ctx context.Context, email string, name string, cognitoSub string, defaultCurrency string) (*entity.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	now := time.Now().UTC()
	user = &entity.User{
		ID:              uuid.NewString(),
		Email:           email,
		Name:            name,
		CognitoSub:      cognitoSub,
		DefaultCurrency: entity.Currency(defaultCurrency),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
