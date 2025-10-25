package dto

type CreateAccountRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=checking savings credit cash"`
	Currency    string  `json:"currency" binding:"required,len=3"`
	Description string  `json:"description"`
	Balance     float64 `json:"balance"`
}

type UpdateAccountRequest struct {
	Name        *string `json:"name"`
	Type        *string `json:"type"`
	Currency    *string `json:"currency"`
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
