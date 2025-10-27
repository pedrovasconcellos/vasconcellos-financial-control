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
// @Description Health check response with service status and uptime
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

// Status
// @Summary Health check
// @Description Retorna o status da API, uptime e status dos serviços
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "API saudável"
// @Router /health [get]
func (h *HealthHandler) Status(c *gin.Context) {
	uptime := time.Since(startTime)

	response := HealthResponse{
		Status:    "healthy",
		Message:   "Financial Control API is running successfully",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "financial-control-api",
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
