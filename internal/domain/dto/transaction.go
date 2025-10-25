package dto

import "time"

type CreateTransactionRequest struct {
	AccountID   string    `json:"accountId" binding:"required"`
	CategoryID  string    `json:"categoryId" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	Currency    string    `json:"currency" binding:"required,len=3"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurredAt" binding:"required"`
	Tags        []string  `json:"tags"`
	Notes       string    `json:"notes"`
}

type UpdateTransactionRequest struct {
	CategoryID  *string  `json:"categoryId"`
	Description *string  `json:"description"`
	Tags        []string `json:"tags"`
	Notes       *string  `json:"notes"`
	Status      *string  `json:"status"`
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	CategoryID  string    `json:"categoryId"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurredAt"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags"`
	Notes       string    `json:"notes"`
	ReceiptURL  *string   `json:"receiptUrl"`
}
