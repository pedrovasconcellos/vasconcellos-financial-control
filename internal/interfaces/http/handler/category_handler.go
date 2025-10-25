package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.categoryUseCase.CreateCategory(c.Request.Context(), user.ID, request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CategoryHandler) List(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := h.categoryUseCase.ListCategories(c.Request.Context(), user.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.categoryUseCase.DeleteCategory(c.Request.Context(), user.ID, c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
