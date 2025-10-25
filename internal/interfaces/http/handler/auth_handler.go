package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Login(c *gin.Context) {
	log := middleware.LoggerFromContext(c)

	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid login request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("authenticating user", zap.String("username", request.Username))
	response, err := h.authUseCase.Login(c.Request.Context(), request)
	if err != nil {
		log.Warn("authentication failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	log.Info("user authenticated")
	c.JSON(http.StatusOK, response)
}
