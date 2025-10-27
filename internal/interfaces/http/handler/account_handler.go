package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type AccountHandler struct {
	accountUseCase *usecase.AccountUseCase
}

func NewAccountHandler(accountUseCase *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{accountUseCase: accountUseCase}
}

// Create
// @Summary Create a new account
// @Description Cria uma nova conta bancária para o usuário autenticado
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAccountRequest true "Dados da conta"
// @Success 201 {object} dto.AccountResponse "Conta criada com sucesso"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized account creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid account create payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("creating account", zap.String("user_id", user.ID), zap.String("name", request.Name))
	response, err := h.accountUseCase.CreateAccount(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to create account", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("account created", zap.String("account_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

// Update
// @Summary Update an account
// @Description Atualiza uma conta existente
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Param request body dto.UpdateAccountRequest true "Dados atualizados"
// @Success 200 {object} dto.AccountResponse "Conta atualizada"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 404 {object} ErrorResponse "Conta não encontrada"
// @Router /accounts/{id} [patch]
func (h *AccountHandler) Update(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized account update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid account update payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountID := c.Param("id")
	log.Info("updating account", zap.String("user_id", user.ID), zap.String("account_id", accountID))
	response, err := h.accountUseCase.UpdateAccount(c.Request.Context(), user.ID, c.Param("id"), request)
	if err != nil {
		log.Error("failed to update account", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// List
// @Summary List accounts
// @Description Lista todas as contas do usuário autenticado
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Número máximo de resultados (default: 100, max: 200)"
// @Param offset query int false "Número de resultados para pular (default: 0)"
// @Success 200 {array} dto.AccountResponse "Lista de contas"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized account list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, offset, err := parsePagination(c.Query("limit"), c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("listing accounts", zap.String("user_id", user.ID), zap.Int64("limit", limit), zap.Int64("offset", offset))
	response, err := h.accountUseCase.ListAccounts(c.Request.Context(), user.ID, limit, offset)
	if err != nil {
		log.Error("failed to list accounts", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete
// @Summary Delete an account
// @Description Remove uma conta existente
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conta"
// @Success 204 "Conta removida com sucesso"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 404 {object} ErrorResponse "Conta não encontrada"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized account delete attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accountID := c.Param("id")
	log.Info("deleting account", zap.String("user_id", user.ID), zap.String("account_id", accountID))
	if err := h.accountUseCase.DeleteAccount(c.Request.Context(), user.ID, accountID); err != nil {
		log.Error("failed to delete account", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("account deleted", zap.String("account_id", accountID))
	c.Status(http.StatusNoContent)
}
