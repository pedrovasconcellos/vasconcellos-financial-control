package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type GoalHandler struct {
	goalUseCase *usecase.GoalUseCase
}

func NewGoalHandler(goalUseCase *usecase.GoalUseCase) *GoalHandler {
	return &GoalHandler{goalUseCase: goalUseCase}
}

func (h *GoalHandler) Create(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateGoalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.goalUseCase.CreateGoal(c.Request.Context(), user.ID, request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *GoalHandler) List(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := h.goalUseCase.ListGoals(c.Request.Context(), user.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GoalHandler) UpdateProgress(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateGoalProgressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.goalUseCase.UpdateProgress(c.Request.Context(), user.ID, c.Param("id"), request.Amount)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
