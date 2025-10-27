package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type GoalHandler struct {
	goalUseCase *usecase.GoalUseCase
}

func NewGoalHandler(goalUseCase *usecase.GoalUseCase) *GoalHandler {
	return &GoalHandler{goalUseCase: goalUseCase}
}

func (h *GoalHandler) Create(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized goal creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateGoalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid goal payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("creating goal", zap.String("user_id", user.ID), zap.String("name", request.Name))
	response, err := h.goalUseCase.CreateGoal(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to create goal", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("goal created", zap.String("goal_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

func (h *GoalHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized goal list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	log.Info("listing goals", zap.String("user_id", user.ID))
	response, err := h.goalUseCase.ListGoals(c.Request.Context(), user.ID)
	if err != nil {
		log.Error("failed to list goals", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GoalHandler) UpdateProgress(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized goal progress attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateGoalProgressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid goal progress payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goalID := c.Param("id")
	log.Info("updating goal progress", zap.String("goal_id", goalID), zap.String("user_id", user.ID), zap.Float64("amount", request.Amount))
	response, err := h.goalUseCase.UpdateProgress(c.Request.Context(), user.ID, c.Param("id"), request.Amount)
	if err != nil {
		log.Error("failed to update goal progress", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
