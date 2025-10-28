package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/src/internal/adapters/http/middleware"
	_ "github.com/vasconcellos/financial-control/src/internal/domain/dto"
	"github.com/vasconcellos/financial-control/src/internal/usecase"
)

type ReportHandler struct {
	reportUseCase *usecase.ReportUseCase
}

func NewReportHandler(reportUseCase *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUseCase: reportUseCase}
}

// Summary
// @Summary Get financial summary
// @Description Gera um resumo financeiro com receitas, despesas e saldo do período
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Data inicial (RFC3339, default: 30 dias atrás)"
// @Param to query string false "Data final (RFC3339, default: hoje)"
// @Success 200 {object} dto.SummaryReportResponse "Resumo financeiro"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Router /reports/summary [get]
func (h *ReportHandler) Summary(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized summary report attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to := parseSummaryRange(c.Query("from"), c.Query("to"))
	log.Info("generating summary report", zap.String("user_id", user.ID), zap.Time("from", from), zap.Time("to", to))
	response, err := h.reportUseCase.GetSummary(c.Request.Context(), user.ID, from, to)
	if err != nil {
		log.Error("failed to generate summary report", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func parseSummaryRange(fromRaw, toRaw string) (time.Time, time.Time) {
	const layout = time.RFC3339
	now := time.Now().UTC()
	if fromRaw == "" {
		fromRaw = now.AddDate(0, -1, 0).Format(layout)
	}
	if toRaw == "" {
		toRaw = now.Format(layout)
	}

	from, err := time.Parse(layout, fromRaw)
	if err != nil {
		from = now.AddDate(0, -1, 0)
	}
	to, err := time.Parse(layout, toRaw)
	if err != nil {
		to = now
	}
	return from, to
}
