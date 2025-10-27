package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/config"
	"github.com/vasconcellos/financial-control/internal/infrastructure/mongodb"
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
	mongoClient  *mongodb.Client
	budgetRepo   *mongodb.BudgetRepository
	lambdaLogger *zap.Logger
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	lambdaLogger = logger
	lambdaLogger.Info("lambda bootstrapped")
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	if mongoClient == nil {
		if err := initDependencies(ctx); err != nil {
			lambdaLogger.Error("failed to initialize dependencies", zap.Error(err))
			return err
		}
	}

	lambdaLogger.Info("processing sqs batch", zap.Int("records", len(event.Records)))
	for _, record := range event.Records {
		var payload transactionEvent
		if err := json.Unmarshal([]byte(record.Body), &payload); err != nil {
			lambdaLogger.Warn("failed to parse message", zap.String("message_id", record.MessageId), zap.Error(err))
			continue
		}

		if err := processTransaction(ctx, payload); err != nil {
			lambdaLogger.Error("failed to process transaction", zap.String("transaction_id", payload.TransactionID), zap.Error(err))
		} else {
			lambdaLogger.Info("transaction processed", zap.String("transaction_id", payload.TransactionID))
		}
	}

	return nil
}

func initDependencies(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lambdaLogger.Info("connecting mongo", zap.String("uri", cfg.Mongo.URI))
	mongoClient, err = mongodb.NewClient(ctx, cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		return err
	}

	budgetRepo = mongodb.NewBudgetRepository(mongoClient)
	lambdaLogger.Info("dependencies initialized")
	return nil
}

func processTransaction(ctx context.Context, payload transactionEvent) error {
	if payload.Type != "expense" {
		lambdaLogger.Debug("ignoring non-expense transaction", zap.String("transaction_id", payload.TransactionID))
		return nil
	}

	lambdaLogger.Info("updating budget spending", zap.String("transaction_id", payload.TransactionID), zap.Float64("amount", payload.Amount))
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
	lambdaLogger.Info("budget spending updated", zap.Int("budgets", len(budgets)))
	return nil
}
