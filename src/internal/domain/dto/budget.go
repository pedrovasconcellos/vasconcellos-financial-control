package dto

import "time"

type CreateBudgetRequest struct {
	CategoryID   string    `json:"categoryId" binding:"required"`
	Amount       float64   `json:"amount" binding:"required"`
	Currency     string    `json:"currency" binding:"required,oneof=USD EUR CHF GBP BRL"`
	Period       string    `json:"period" binding:"required,oneof=monthly quarterly yearly"`
	PeriodStart  time.Time `json:"periodStart" binding:"required"`
	PeriodEnd    time.Time `json:"periodEnd" binding:"required"`
	AlertPercent float64   `json:"alertPercent" binding:"required"`
}

type BudgetResponse struct {
	ID           string    `json:"id"`
	CategoryID   string    `json:"categoryId"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Period       string    `json:"period"`
	PeriodStart  time.Time `json:"periodStart"`
	PeriodEnd    time.Time `json:"periodEnd"`
	Spent        float64   `json:"spent"`
	AlertPercent float64   `json:"alertPercent"`
}
