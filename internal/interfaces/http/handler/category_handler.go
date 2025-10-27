package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

// Create
// @Summary Create a new category
// @Description Cria uma nova categoria de transação (income ou expense)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Dados da categoria"
// @Success 201 {object} dto.CategoryResponse "Categoria criada"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized category creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid category payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("creating category", zap.String("user_id", user.ID), zap.String("name", request.Name))
	response, err := h.categoryUseCase.CreateCategory(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to create category", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("category created", zap.String("category_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

// List
// @Summary List categories
// @Description Lista todas as categorias do usuário
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryResponse "Lista de categorias"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized category list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	log.Info("listing categories", zap.String("user_id", user.ID))
	response, err := h.categoryUseCase.ListCategories(c.Request.Context(), user.ID)
	if err != nil {
		log.Error("failed to list categories", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete
// @Summary Delete a category
// @Description Remove uma categoria
// @Tags categories
// @Security BearerAuth
// @Param id path string true "ID da categoria"
// @Success 204 "Categoria removida"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 404 {object} ErrorResponse "Categoria não encontrada"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized category delete attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryID := c.Param("id")
	log.Info("deleting category", zap.String("user_id", user.ID), zap.String("category_id", categoryID))
	if err := h.categoryUseCase.DeleteCategory(c.Request.Context(), user.ID, categoryID); err != nil {
		log.Error("failed to delete category", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("category deleted", zap.String("category_id", categoryID))
	c.Status(http.StatusNoContent)
}
