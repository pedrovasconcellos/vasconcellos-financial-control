package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type BudgetHandler struct {
	budgetUseCase *usecase.BudgetUseCase
}

func NewBudgetHandler(budgetUseCase *usecase.BudgetUseCase) *BudgetHandler {
	return &BudgetHandler{budgetUseCase: budgetUseCase}
}

// Create
// @Summary Create a budget
// @Description Cria um orçamento para controlar gastos por categoria
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateBudgetRequest true "Dados do orçamento"
// @Success 201 {object} dto.BudgetResponse "Orçamento criado"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Router /budgets [post]
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

// List
// @Summary List budgets
// @Description Lista todos os orçamentos do usuário com status atualizado
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.BudgetResponse "Lista de orçamentos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Router /budgets [get]
func (h *BudgetHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized budget list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, offset, err := parsePagination(c.Query("limit"), c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("listing budgets", zap.String("user_id", user.ID), zap.Int64("limit", limit), zap.Int64("offset", offset))
	response, err := h.budgetUseCase.ListBudgets(c.Request.Context(), user.ID, limit, offset)
	if err != nil {
		log.Error("failed to list budgets", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
