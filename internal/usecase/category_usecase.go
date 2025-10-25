package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/entity"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, userID string, request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	now := time.Now().UTC()
	category := &entity.Category{
		ID:          uuid.NewString(),
		UserID:      userID,
		Name:        request.Name,
		Type:        entity.CategoryType(request.Type),
		Description: request.Description,
		ParentID:    request.ParentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Type:        string(category.Type),
		Description: category.Description,
		ParentID:    category.ParentID,
	}, nil
}

func (uc *CategoryUseCase) ListCategories(ctx context.Context, userID string) ([]*dto.CategoryResponse, error) {
	categories, err := uc.categoryRepo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Type:        string(category.Type),
			Description: category.Description,
			ParentID:    category.ParentID,
		})
	}
	return response, nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, userID string, categoryID string) error {
	return uc.categoryRepo.Delete(ctx, categoryID, userID)
}
