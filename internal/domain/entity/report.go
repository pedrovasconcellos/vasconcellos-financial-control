package entity

type SummaryReport struct {
	TotalIncome        float64            `bson:"total_income"`
	TotalExpense       float64            `bson:"total_expense"`
	NetBalance         float64            `bson:"net_balance"`
	SpendingByCategory map[string]float64 `bson:"spending_by_category"`
	BudgetUsage        map[string]float64 `bson:"budget_usage"`
	GoalProgress       map[string]float64 `bson:"goal_progress"`
}
