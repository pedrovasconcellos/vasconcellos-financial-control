package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "github.com/vasconcellos/financial-control/src/internal/domain/errors"
)

// ErrorResponse representa uma resposta de erro padrão da API
// @Description Standard error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

func respondError(c *gin.Context, err error) {
	switch err {
	case nil:
		c.Status(http.StatusOK)
	case domainErrors.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domainErrors.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case domainErrors.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case domainErrors.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
