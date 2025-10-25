package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type BudgetHandler struct {
	budgetUseCase *usecase.BudgetUseCase
}

func NewBudgetHandler(budgetUseCase *usecase.BudgetUseCase) *BudgetHandler {
	return &BudgetHandler{budgetUseCase: budgetUseCase}
}

func (h *BudgetHandler) Create(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.budgetUseCase.CreateBudget(c.Request.Context(), user.ID, request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *BudgetHandler) List(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := h.budgetUseCase.ListBudgets(c.Request.Context(), user.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
