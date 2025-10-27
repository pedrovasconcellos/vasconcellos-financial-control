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

func (h *AccountHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized account list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	log.Info("listing accounts", zap.String("user_id", user.ID))
	response, err := h.accountUseCase.ListAccounts(c.Request.Context(), user.ID)
	if err != nil {
		log.Error("failed to list accounts", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

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
