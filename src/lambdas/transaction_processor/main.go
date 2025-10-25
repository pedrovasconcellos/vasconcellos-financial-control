package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/vasconcellos/finance-control/internal/config"
	"github.com/vasconcellos/finance-control/internal/infrastructure/mongodb"
)

type transactionEvent struct {
	TransactionID string    `json:"transactionId"`
	UserID        string    `json:"userId"`
	AccountID     string    `json:"accountId"`
	CategoryID    string    `json:"categoryId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	OccurredAt    time.Time `json:"occurredAt"`
	Type          string    `json:"type"`
}

var (
	mongoClient *mongodb.Client
	budgetRepo  *mongodb.BudgetRepository
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	if mongoClient == nil {
		if err := initDependencies(ctx); err != nil {
			return err
		}
	}

	for _, record := range event.Records {
		var payload transactionEvent
		if err := json.Unmarshal([]byte(record.Body), &payload); err != nil {
			log.Printf("failed to parse message: %v", err)
			continue
		}

		if err := processTransaction(ctx, payload); err != nil {
			log.Printf("failed to process transaction %s: %v", payload.TransactionID, err)
		}
	}

	return nil
}

func initDependencies(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	mongoClient, err = mongodb.NewClient(ctx, cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		return err
	}

	budgetRepo = mongodb.NewBudgetRepository(mongoClient)
	return nil
}

func processTransaction(ctx context.Context, payload transactionEvent) error {
	if payload.Type != "expense" {
		return nil
	}

	budgets, err := budgetRepo.FindActiveByCategory(ctx, payload.UserID, payload.CategoryID, payload.OccurredAt)
	if err != nil {
		return err
	}

	for _, budget := range budgets {
		newSpent := budget.Spent + payload.Amount
		if err := budgetRepo.UpdateSpent(ctx, budget.ID, budget.UserID, newSpent); err != nil {
			return err
		}
	}
	return nil
}
