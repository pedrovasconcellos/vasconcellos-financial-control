package entity

import "time"

type BudgetPeriod string

const (
	BudgetPeriodMonthly   BudgetPeriod = "monthly"
	BudgetPeriodQuarterly BudgetPeriod = "quarterly"
	BudgetPeriodYearly    BudgetPeriod = "yearly"
)

type Budget struct {
	ID           string       `bson:"_id,omitempty"`
	UserID       string       `bson:"user_id"`
	CategoryID   string       `bson:"category_id"`
	Amount       float64      `bson:"amount"`
	Currency     string       `bson:"currency"`
	Period       BudgetPeriod `bson:"period"`
	PeriodStart  time.Time    `bson:"period_start"`
	PeriodEnd    time.Time    `bson:"period_end"`
	Spent        float64      `bson:"spent"`
	CreatedAt    time.Time    `bson:"created_at"`
	UpdatedAt    time.Time    `bson:"updated_at"`
	AlertPercent float64      `bson:"alert_percent"`
}
