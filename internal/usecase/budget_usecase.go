package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/entity"
	"github.com/vasconcellos/finance-control/internal/domain/errors"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type BudgetUseCase struct {
	budgetRepo repository.BudgetRepository
}

func NewBudgetUseCase(budgetRepo repository.BudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{
		budgetRepo: budgetRepo,
	}
}

func (uc *BudgetUseCase) CreateBudget(ctx context.Context, userID string, request dto.CreateBudgetRequest) (*dto.BudgetResponse, error) {
	budget := &entity.Budget{
		ID:           uuid.NewString(),
		UserID:       userID,
		CategoryID:   request.CategoryID,
		Amount:       request.Amount,
		Currency:     entity.Currency(request.Currency),
		Period:       entity.BudgetPeriod(request.Period),
		PeriodStart:  request.PeriodStart,
		PeriodEnd:    request.PeriodEnd,
		AlertPercent: request.AlertPercent,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := uc.budgetRepo.Create(ctx, budget); err != nil {
		return nil, err
	}

	return &dto.BudgetResponse{
		ID:           budget.ID,
		CategoryID:   budget.CategoryID,
		Amount:       budget.Amount,
		Currency:     budget.Currency.String(),
		Period:       string(budget.Period),
		PeriodStart:  budget.PeriodStart,
		PeriodEnd:    budget.PeriodEnd,
		Spent:        budget.Spent,
		AlertPercent: budget.AlertPercent,
	}, nil
}

func (uc *BudgetUseCase) ListBudgets(ctx context.Context, userID string) ([]*dto.BudgetResponse, error) {
	budgets, err := uc.budgetRepo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.BudgetResponse, 0, len(budgets))
	for _, budget := range budgets {
		response = append(response, &dto.BudgetResponse{
			ID:           budget.ID,
			CategoryID:   budget.CategoryID,
			Amount:       budget.Amount,
			Currency:     budget.Currency.String(),
			Period:       string(budget.Period),
			PeriodStart:  budget.PeriodStart,
			PeriodEnd:    budget.PeriodEnd,
			Spent:        budget.Spent,
			AlertPercent: budget.AlertPercent,
		})
	}

	return response, nil
}

func (uc *BudgetUseCase) UpdateSpent(ctx context.Context, userID string, budgetID string, spent float64) error {
	budget, err := uc.budgetRepo.GetByID(ctx, budgetID, userID)
	if err != nil {
		return err
	}
	if budget == nil {
		return errors.ErrNotFound
	}

	return uc.budgetRepo.UpdateSpent(ctx, budgetID, userID, spent)
}
