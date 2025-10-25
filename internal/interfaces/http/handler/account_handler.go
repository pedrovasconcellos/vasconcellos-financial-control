package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type AccountHandler struct {
	accountUseCase *usecase.AccountUseCase
}

func NewAccountHandler(accountUseCase *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{accountUseCase: accountUseCase}
}

func (h *AccountHandler) Create(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.accountUseCase.CreateAccount(c.Request.Context(), user.ID, request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AccountHandler) Update(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.accountUseCase.UpdateAccount(c.Request.Context(), user.ID, c.Param("id"), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AccountHandler) List(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := h.accountUseCase.ListAccounts(c.Request.Context(), user.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.accountUseCase.DeleteAccount(c.Request.Context(), user.ID, c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
