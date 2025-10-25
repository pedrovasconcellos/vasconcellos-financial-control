package entity

import "time"

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID          string       `bson:"_id,omitempty"`
	UserID      string       `bson:"user_id"`
	Name        string       `bson:"name"`
	Type        CategoryType `bson:"type"`
	ParentID    *string      `bson:"parent_id,omitempty"`
	Description string       `bson:"description"`
	CreatedAt   time.Time    `bson:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at"`
}
