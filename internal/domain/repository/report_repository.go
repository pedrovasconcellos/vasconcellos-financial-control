package repository

import (
	"context"
	"time"

	"github.com/vasconcellos/finance-control/internal/domain/entity"
)

type ReportRepository interface {
	AggregateSummary(ctx context.Context, userID string, from time.Time, to time.Time) (*entity.SummaryReport, error)
}
