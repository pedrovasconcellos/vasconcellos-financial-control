package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized budget creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid budget payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("creating budget", zap.String("user_id", user.ID), zap.String("category_id", request.CategoryID))
	response, err := h.budgetUseCase.CreateBudget(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to create budget", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("budget created", zap.String("budget_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

func (h *BudgetHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized budget list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	log.Info("listing budgets", zap.String("user_id", user.ID))
	response, err := h.budgetUseCase.ListBudgets(c.Request.Context(), user.ID)
	if err != nil {
		log.Error("failed to list budgets", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
