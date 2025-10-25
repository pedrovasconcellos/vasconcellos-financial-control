package entity

import "time"

type AccountType string

const (
	AccountTypeChecking AccountType = "checking"
	AccountTypeSavings  AccountType = "savings"
	AccountTypeCredit   AccountType = "credit"
	AccountTypeCash     AccountType = "cash"
)

type Account struct {
	ID          string      `bson:"_id,omitempty"`
	UserID      string      `bson:"user_id"`
	Name        string      `bson:"name"`
	Type        AccountType `bson:"type"`
	Currency    string      `bson:"currency"`
	Balance     float64     `bson:"balance"`
	Description string      `bson:"description"`
	CreatedAt   time.Time   `bson:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at"`
}
