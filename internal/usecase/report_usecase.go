package usecase

import (
	"context"
	"time"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type ReportUseCase struct {
	reportRepo repository.ReportRepository
}

func NewReportUseCase(reportRepo repository.ReportRepository) *ReportUseCase {
	return &ReportUseCase{reportRepo: reportRepo}
}

func (uc *ReportUseCase) GetSummary(ctx context.Context, userID string, from, to time.Time) (*dto.SummaryReportResponse, error) {
	report, err := uc.reportRepo.AggregateSummary(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	return &dto.SummaryReportResponse{
		TotalIncome:        report.TotalIncome,
		TotalExpense:       report.TotalExpense,
		NetBalance:         report.NetBalance,
		SpendingByCategory: report.SpendingByCategory,
		BudgetUsage:        report.BudgetUsage,
		GoalProgress:       report.GoalProgress,
	}, nil
}
