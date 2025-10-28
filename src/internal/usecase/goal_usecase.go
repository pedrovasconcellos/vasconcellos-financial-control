package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/financial-control/src/internal/domain/dto"
	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
	"github.com/vasconcellos/financial-control/src/internal/domain/errors"
	"github.com/vasconcellos/financial-control/src/internal/domain/repository"
)

type GoalUseCase struct {
	goalRepo repository.GoalRepository
}

func NewGoalUseCase(goalRepo repository.GoalRepository) *GoalUseCase {
	return &GoalUseCase{
		goalRepo: goalRepo,
	}
}

func (uc *GoalUseCase) CreateGoal(ctx context.Context, userID string, request dto.CreateGoalRequest) (*dto.GoalResponse, error) {
	now := time.Now().UTC()
	goal := &entity.Goal{
		ID:            uuid.NewString(),
		UserID:        userID,
		Name:          request.Name,
		TargetAmount:  request.TargetAmount,
		CurrentAmount: 0,
		Currency:      entity.Currency(request.Currency),
		Deadline:      request.Deadline,
		Description:   request.Description,
		Status:        entity.GoalStatusActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.goalRepo.Create(ctx, goal); err != nil {
		return nil, err
	}

	return &dto.GoalResponse{
		ID:            goal.ID,
		Name:          goal.Name,
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		Currency:      goal.Currency.String(),
		Deadline:      goal.Deadline,
		Status:        string(goal.Status),
		Description:   goal.Description,
	}, nil
}

func (uc *GoalUseCase) ListGoals(ctx context.Context, userID string, limit int64, offset int64) ([]*dto.GoalResponse, error) {
	goals, err := uc.goalRepo.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.GoalResponse, 0, len(goals))
	for _, goal := range goals {
		response = append(response, &dto.GoalResponse{
			ID:            goal.ID,
			Name:          goal.Name,
			TargetAmount:  goal.TargetAmount,
			CurrentAmount: goal.CurrentAmount,
			Currency:      goal.Currency.String(),
			Deadline:      goal.Deadline,
			Status:        string(goal.Status),
			Description:   goal.Description,
		})
	}

	return response, nil
}

func (uc *GoalUseCase) UpdateProgress(ctx context.Context, userID string, goalID string, amount float64) (*dto.GoalResponse, error) {
	goal, err := uc.goalRepo.GetByID(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}
	if goal == nil {
		return nil, errors.ErrNotFound
	}

	goal.CurrentAmount += amount
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = entity.GoalStatusCompleted
	}
	goal.UpdatedAt = time.Now().UTC()

	if err := uc.goalRepo.Update(ctx, goal); err != nil {
		return nil, err
	}

	return &dto.GoalResponse{
		ID:            goal.ID,
		Name:          goal.Name,
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		Currency:      goal.Currency.String(),
		Deadline:      goal.Deadline,
		Status:        string(goal.Status),
		Description:   goal.Description,
	}, nil
}
