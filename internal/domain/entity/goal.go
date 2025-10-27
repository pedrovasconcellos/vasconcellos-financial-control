package entity

import "time"

type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusOnHold    GoalStatus = "on_hold"
)

type Goal struct {
	ID            string     `bson:"_id"`
	UserID        string     `bson:"user_id"`
	Name          string     `bson:"name"`
	TargetAmount  float64    `bson:"target_amount"`
	CurrentAmount float64    `bson:"current_amount"`
	Currency      Currency   `bson:"currency"`
	Deadline      time.Time  `bson:"deadline"`
	Status        GoalStatus `bson:"status"`
	Description   string     `bson:"description"`
	CreatedAt     time.Time  `bson:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at"`
}
