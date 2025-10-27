package dto

import "time"

type CreateGoalRequest struct {
	Name         string    `json:"name" binding:"required"`
	TargetAmount float64   `json:"targetAmount" binding:"required"`
	Currency     string    `json:"currency" binding:"required,oneof=USD EUR CHF GBP BRL"`
	Deadline     time.Time `json:"deadline" binding:"required"`
	Description  string    `json:"description"`
}

type UpdateGoalProgressRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type GoalResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"targetAmount"`
	CurrentAmount float64   `json:"currentAmount"`
	Currency      string    `json:"currency"`
	Deadline      time.Time `json:"deadline"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
}
