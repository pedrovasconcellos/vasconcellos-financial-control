package main

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/src/internal/config"
	"github.com/vasconcellos/financial-control/src/internal/infrastructure/mongodb"
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
	mongoClient   *mongodb.Client
	budgetRepo    *mongodb.BudgetRepository
	processedRepo *mongodb.ProcessedTransactionRepository
	sqsClient     *sqs.Client
	queueURL      string
	awsCfg        aws.Config
	appConfig     *config.Config
	lambdaLogger  *zap.Logger
)

const (
	localModeEnv          = "LAMBDA_LOCAL"
	defaultWaitTimeSecond = 10
	maxSQSBatchSize       = 10
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	lambdaLogger = logger
	lambdaLogger.Info("lambda bootstrapped")

	if strings.EqualFold(os.Getenv(localModeEnv), "true") {
		ctx := context.Background()
		if err := initDependencies(ctx); err != nil {
			lambdaLogger.Fatal("failed to initialize dependencies", zap.Error(err))
		}
		if err := ensureSQSClient(ctx); err != nil {
			lambdaLogger.Fatal("failed to initialize sqs client", zap.Error(err))
		}
		lambdaLogger.Info("starting local lambda worker", zap.String("queue_url", queueURL))
		startLocalWorker(ctx)
		return
	}

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
	if appConfig != nil {
		return nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	appConfig = cfg

	lambdaLogger.Info("connecting mongo", zap.String("uri", cfg.Mongo.URI))
	mongoClient, err = mongodb.NewClient(ctx, cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		return err
	}

	budgetRepo = mongodb.NewBudgetRepository(mongoClient)
	processedRepo = mongodb.NewProcessedTransactionRepository(mongoClient)

	awsCfg, err = buildAWSConfig(ctx, cfg)
	if err != nil {
		return err
	}
	queueURL = cfg.AWS.SQS.QueueURL
	lambdaLogger.Info("dependencies initialized")
	return nil
}

func buildAWSConfig(ctx context.Context, cfg *config.Config) (aws.Config, error) {
	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.AWS.Region),
	}
	if cfg.AWS.Endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, args ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: cfg.AWS.Endpoint, SigningRegion: cfg.AWS.Region}, nil
		})
		options = append(options, awsconfig.WithEndpointResolverWithOptions(resolver))
	}
	if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		provider := credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.SessionToken)
		options = append(options, awsconfig.WithCredentialsProvider(provider))
	}

	return awsconfig.LoadDefaultConfig(ctx, options...)
}

func ensureSQSClient(ctx context.Context) error {
	if sqsClient != nil {
		return nil
	}
	sqsClient = sqs.NewFromConfig(awsCfg)
	if queueURL == "" {
		return errors.New("sqs queue url not configured")
	}
	return nil
}

func processTransaction(ctx context.Context, payload transactionEvent) error {
	inserted, err := processedRepo.MarkProcessed(ctx, payload.TransactionID, payload.UserID, payload.Type, time.Now().UTC())
	if err != nil {
		return err
	}
	if !inserted {
		lambdaLogger.Debug("transaction already processed", zap.String("transaction_id", payload.TransactionID))
		return nil
	}

	amount := math.Abs(payload.Amount)
	var delta float64
	switch payload.Type {
	case "expense":
		delta = amount
	case "income":
		delta = -amount
	default:
		lambdaLogger.Warn("unknown transaction type", zap.String("transaction_id", payload.TransactionID), zap.String("type", payload.Type))
		return nil
	}

	if delta == 0 {
		lambdaLogger.Debug("ignoring zero-impact transaction", zap.String("transaction_id", payload.TransactionID))
		return nil
	}

	lambdaLogger.Info("updating budget spending", zap.String("transaction_id", payload.TransactionID), zap.Float64("delta", delta))
	budgets, err := budgetRepo.FindActiveByCategory(ctx, payload.UserID, payload.CategoryID, payload.OccurredAt)
	if err != nil {
		if removeErr := processedRepo.Remove(ctx, payload.TransactionID); removeErr != nil {
			lambdaLogger.Warn("failed to rollback processed marker", zap.String("transaction_id", payload.TransactionID), zap.Error(removeErr))
		}
		return err
	}

	for _, budget := range budgets {
		newSpent := budget.Spent + delta
		if newSpent < 0 {
			newSpent = 0
		}
		if err := budgetRepo.UpdateSpent(ctx, budget.ID, budget.UserID, newSpent); err != nil {
			if removeErr := processedRepo.Remove(ctx, payload.TransactionID); removeErr != nil {
				lambdaLogger.Warn("failed to rollback processed marker", zap.String("transaction_id", payload.TransactionID), zap.Error(removeErr))
			}
			return err
		}
	}
	lambdaLogger.Info("budget spending updated", zap.Int("budgets", len(budgets)))
	return nil
}

func startLocalWorker(ctx context.Context) {
	for {
		output, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: maxSQSBatchSize,
			WaitTimeSeconds:     defaultWaitTimeSecond,
		})
		if err != nil {
			lambdaLogger.Error("failed to receive messages", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}

		if len(output.Messages) == 0 {
			continue
		}

		for _, message := range output.Messages {
			body := aws.ToString(message.Body)
			if body == "" {
				deleteMessage(ctx, message)
				continue
			}
			var payload transactionEvent
			if err := json.Unmarshal([]byte(body), &payload); err != nil {
				lambdaLogger.Warn("failed to decode message", zap.Error(err))
				deleteMessage(ctx, message)
				continue
			}

			if err := processTransaction(ctx, payload); err != nil {
				lambdaLogger.Error("failed to process transaction", zap.String("transaction_id", payload.TransactionID), zap.Error(err))
				continue
			}

			deleteMessage(ctx, message)
		}
	}
}

func deleteMessage(ctx context.Context, message types.Message) {
	if message.ReceiptHandle == nil {
		return
	}
	_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		lambdaLogger.Warn("failed to delete message", zap.Error(err))
	}
}
