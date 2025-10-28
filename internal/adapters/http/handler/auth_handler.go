package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/adapters/http/middleware"
	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// Login
// @Summary Authenticate user
// @Description Autentica um usuário e retorna tokens de acesso
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciais de login"
// @Success 200 {object} dto.LoginResponse "Tokens de autenticação"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Credenciais inválidas"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /auth/login [post]
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
		status := http.StatusUnauthorized
		message := "authentication failed"
		if err == errors.ErrInvalidInput {
			status = http.StatusBadRequest
			message = "invalid credentials"
		}
		log.Warn("authentication failed", zap.Error(err))
		c.JSON(status, gin.H{"error": message})
		return
	}

	log.Info("user authenticated")
	c.JSON(http.StatusOK, response)
}
