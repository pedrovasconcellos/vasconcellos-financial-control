package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type ReportHandler struct {
	reportUseCase *usecase.ReportUseCase
}

func NewReportHandler(reportUseCase *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUseCase: reportUseCase}
}

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
