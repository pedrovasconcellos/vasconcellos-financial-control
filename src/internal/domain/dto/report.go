package dto

type SummaryReportResponse struct {
	TotalIncome        float64            `json:"totalIncome"`
	TotalExpense       float64            `json:"totalExpense"`
	NetBalance         float64            `json:"netBalance"`
	SpendingByCategory map[string]float64 `json:"spendingByCategory"`
	BudgetUsage        map[string]float64 `json:"budgetUsage"`
	GoalProgress       map[string]float64 `json:"goalProgress"`
}
