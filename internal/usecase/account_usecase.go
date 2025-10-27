package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/domain/entity"
	"github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
)

type AccountUseCase struct {
	accountRepo repository.AccountRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
	}
}

func (uc *AccountUseCase) CreateAccount(ctx context.Context, userID string, request dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	now := time.Now().UTC()
	account := &entity.Account{
		ID:          uuid.NewString(),
		UserID:      userID,
		Name:        request.Name,
		Type:        entity.AccountType(request.Type),
		Currency:    entity.Currency(request.Currency),
		Description: request.Description,
		Balance:     request.Balance,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return &dto.AccountResponse{
		ID:          account.ID,
		Name:        account.Name,
		Type:        string(account.Type),
		Currency:    account.Currency.String(),
		Description: account.Description,
		Balance:     account.Balance,
	}, nil
}

func (uc *AccountUseCase) UpdateAccount(ctx context.Context, userID string, accountID string, request dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	account, err := uc.accountRepo.GetByID(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.ErrNotFound
	}

	if request.Name != nil {
		account.Name = *request.Name
	}
	if request.Type != nil {
		account.Type = entity.AccountType(*request.Type)
	}
	if request.Currency != nil {
		account.Currency = entity.Currency(*request.Currency)
	}
	if request.Description != nil {
		account.Description = *request.Description
	}
	account.UpdatedAt = time.Now().UTC()

	if err := uc.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return &dto.AccountResponse{
		ID:          account.ID,
		Name:        account.Name,
		Type:        string(account.Type),
		Currency:    account.Currency.String(),
		Description: account.Description,
		Balance:     account.Balance,
	}, nil
}

func (uc *AccountUseCase) ListAccounts(ctx context.Context, userID string, limit int64, offset int64) ([]*dto.AccountResponse, error) {
	accounts, err := uc.accountRepo.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		response = append(response, &dto.AccountResponse{
			ID:          account.ID,
			Name:        account.Name,
			Type:        string(account.Type),
			Currency:    account.Currency.String(),
			Description: account.Description,
			Balance:     account.Balance,
		})
	}

	return response, nil
}

func (uc *AccountUseCase) DeleteAccount(ctx context.Context, userID string, accountID string) error {
	return uc.accountRepo.Delete(ctx, accountID, userID)
}

func (uc *AccountUseCase) AdjustAccountBalance(ctx context.Context, userID string, accountID string, amount float64) error {
	return uc.accountRepo.AdjustBalance(ctx, accountID, userID, amount)
}
