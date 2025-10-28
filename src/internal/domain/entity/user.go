package entity

import "time"

type User struct {
	ID              string    `bson:"_id"`
	Email           string    `bson:"email"`
	Name            string    `bson:"name"`
	DefaultCurrency Currency  `bson:"default_currency"`
	CognitoSub      string    `bson:"cognito_sub"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}
