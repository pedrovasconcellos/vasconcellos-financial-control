package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
)

type ReportRepository struct {
	client *Client
}

var _ repository.ReportRepository = (*ReportRepository)(nil)

func NewReportRepository(client *Client) *ReportRepository {
	return &ReportRepository{client: client}
}

func (r *ReportRepository) AggregateSummary(ctx context.Context, userID string, from time.Time, to time.Time) (*entity.SummaryReport, error) {
	transactions, err := r.fetchTransactions(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	categories, err := r.fetchCategories(ctx, userID)
	if err != nil {
		return nil, err
	}
	budgets, err := r.fetchBudgets(ctx, userID)
	if err != nil {
		return nil, err
	}
	goals, err := r.fetchGoals(ctx, userID)
	if err != nil {
		return nil, err
	}

	report := &entity.SummaryReport{
		SpendingByCategory: map[string]float64{},
		BudgetUsage:        map[string]float64{},
		GoalProgress:       map[string]float64{},
	}

	for _, transaction := range transactions {
		category := categories[transaction.CategoryID]
		if category.Type == entity.CategoryTypeIncome {
			report.TotalIncome += transaction.Amount
		} else {
			report.TotalExpense += transaction.Amount
			name := category.Name
			if name == "" {
				name = transaction.CategoryID
			}
			report.SpendingByCategory[name] += transaction.Amount
		}
	}
	report.NetBalance = report.TotalIncome - report.TotalExpense

	for _, budget := range budgets {
		if budget.Amount == 0 {
			continue
		}
		category := categories[budget.CategoryID]
		name := category.Name
		if name == "" {
			name = budget.CategoryID
		}
		report.BudgetUsage[name] = (budget.Spent / budget.Amount) * 100
	}

	for _, goal := range goals {
		if goal.TargetAmount == 0 {
			continue
		}
		name := goal.Name
		if name == "" {
			name = goal.ID
		}
		report.GoalProgress[name] = (goal.CurrentAmount / goal.TargetAmount) * 100
	}

	return report, nil
}

func (r *ReportRepository) fetchTransactions(ctx context.Context, userID string, from time.Time, to time.Time) ([]*entity.Transaction, error) {
	col := r.client.Collection("transactions")
	cursor, err := col.Find(ctx, bson.M{
		"user_id": userID,
		"occurred_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*entity.Transaction
	for cursor.Next(ctx) {
		var transaction entity.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (r *ReportRepository) fetchCategories(ctx context.Context, userID string) (map[string]entity.Category, error) {
	col := r.client.Collection("categories")
	cursor, err := col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]entity.Category)
	for cursor.Next(ctx) {
		var category entity.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		result[category.ID] = category
	}
	return result, nil
}

func (r *ReportRepository) fetchBudgets(ctx context.Context, userID string) ([]*entity.Budget, error) {
	col := r.client.Collection("budgets")
	cursor, err := col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []*entity.Budget
	for cursor.Next(ctx) {
		var budget entity.Budget
		if err := cursor.Decode(&budget); err != nil {
			return nil, err
		}
		budgets = append(budgets, &budget)
	}
	return budgets, nil
}

func (r *ReportRepository) fetchGoals(ctx context.Context, userID string) ([]*entity.Goal, error) {
	col := r.client.Collection("goals")
	cursor, err := col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var goals []*entity.Goal
	for cursor.Next(ctx) {
		var goal entity.Goal
		if err := cursor.Decode(&goal); err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}
