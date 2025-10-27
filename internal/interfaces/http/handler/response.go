package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
)

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
