package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler responde status básico da aplicação
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
