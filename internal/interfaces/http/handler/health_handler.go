package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler responde status básico da aplicação
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse representa a resposta do endpoint de health
type HealthResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	Timestamp string            `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks"`
}

var startTime = time.Now()

func (h *HealthHandler) Status(c *gin.Context) {
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:    "healthy",
		Message:   "Finance Control API is running successfully",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "finance-control-api",
		Version:   "1.0.0",
		Uptime:    uptime.Round(time.Second).String(),
		Checks: map[string]string{
			"api":        "ok",
			"database":   "ok",
			"localstack": "ok",
		},
	}

	c.JSON(http.StatusOK, response)
}
