package dto

type CreateAccountRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=checking savings credit cash"`
	Currency    string  `json:"currency" binding:"required,oneof=USD EUR CHF GBP BRL"`
	Description string  `json:"description"`
	Balance     float64 `json:"balance"`
}

type UpdateAccountRequest struct {
	Name        *string `json:"name"`
	Type        *string `json:"type"`
	Currency    *string `json:"currency" binding:"omitempty,oneof=USD EUR CHF GBP BRL"`
	Description *string `json:"description"`
}

type AccountResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Balance     float64 `json:"balance"`
}
