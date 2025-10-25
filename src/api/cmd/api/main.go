package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/zap"

	"github.com/vasconcellos/finance-control/internal/config"
	"github.com/vasconcellos/finance-control/internal/domain/port"
	"github.com/vasconcellos/finance-control/internal/infrastructure/auth"
	awsS3 "github.com/vasconcellos/finance-control/internal/infrastructure/aws/s3"
	awsSQS "github.com/vasconcellos/finance-control/internal/infrastructure/aws/sqs"
	"github.com/vasconcellos/finance-control/internal/infrastructure/logger"
	"github.com/vasconcellos/finance-control/internal/infrastructure/mongodb"
	interfacesHTTP "github.com/vasconcellos/finance-control/internal/interfaces/http"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/handler"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logr, err := logger.New(cfg.App.Environment)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logr.Sync()

	mongoClient, err := mongodb.NewClient(ctx, cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		logr.Fatal("failed to connect mongo", zap.Error(err))
	}
	defer mongoClient.Close(context.Background())

	userRepo, err := mongodb.NewUserRepository(mongoClient)
	if err != nil {
		logr.Fatal("failed to init user repo", zap.Error(err))
	}
	accountRepo, err := mongodb.NewAccountRepository(mongoClient)
	if err != nil {
		logr.Fatal("failed to init account repo", zap.Error(err))
	}
	categoryRepo := mongodb.NewCategoryRepository(mongoClient)
	transactionRepo, err := mongodb.NewTransactionRepository(mongoClient)
	if err != nil {
		logr.Fatal("failed to init transaction repo", zap.Error(err))
	}
	budgetRepo := mongodb.NewBudgetRepository(mongoClient)
	goalRepo := mongodb.NewGoalRepository(mongoClient)
	reportRepo := mongodb.NewReportRepository(mongoClient)

	awsCfg, err := buildAWSConfig(ctx, cfg)
	if err != nil {
		logr.Warn("AWS config fallback", zap.Error(err))
	}

	var storage port.ObjectStorage
	if cfg.Storage.ReceiptBucket != "" && err == nil {
		storage = awsS3.NewStorage(awsCfg, cfg.Storage.ReceiptBucket)
	}

	var queuePublisher port.QueuePublisher
	if cfg.AWS.SQS.QueueURL != "" && err == nil {
		queuePublisher = awsSQS.NewPublisher(awsCfg, cfg.AWS.SQS.QueueURL)
	}

	authProvider, authErr := buildAuthProvider(ctx, cfg, err == nil)
	if authErr != nil {
		logr.Fatal("failed to init auth provider", zap.Error(authErr))
	}

	authUseCase := usecase.NewAuthUseCase(authProvider)
	userUseCase := usecase.NewUserUseCase(userRepo)
	accountUseCase := usecase.NewAccountUseCase(accountRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, accountRepo, categoryRepo, queuePublisher, storage, cfg.Queue.TransactionQueue)
	budgetUseCase := usecase.NewBudgetUseCase(budgetRepo)
	goalUseCase := usecase.NewGoalUseCase(goalRepo)
	reportUseCase := usecase.NewReportUseCase(reportRepo)

	authHandler := handler.NewAuthHandler(authUseCase)
	accountHandler := handler.NewAccountHandler(accountUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	budgetHandler := handler.NewBudgetHandler(budgetUseCase)
	goalHandler := handler.NewGoalHandler(goalUseCase)
	reportHandler := handler.NewReportHandler(reportUseCase)

	authMiddleware := middleware.NewAuthMiddleware(authUseCase, userUseCase)
	router := interfacesHTTP.NewRouter(interfacesHTTP.RouterParams{
		AuthHandler:        authHandler,
		AccountHandler:     accountHandler,
		CategoryHandler:    categoryHandler,
		TransactionHandler: transactionHandler,
		BudgetHandler:      budgetHandler,
		GoalHandler:        goalHandler,
		ReportHandler:      reportHandler,
		AuthMiddleware:     authMiddleware,
		AllowedOrigins:     cfg.Security.AllowedOrigins,
		Logger:            logr,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logr.Info("server started", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logr.Fatal("server failure", zap.Error(err))
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logr.Info("shutdown signal received")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logr.Error("server shutdown error", zap.Error(err))
	}
	logr.Info("server stopped")
}

func buildAWSConfig(ctx context.Context, cfg *config.Config) (aws.Config, error) {
	options := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.AWS.Region),
	}
	if cfg.AWS.Endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, args ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: cfg.AWS.Endpoint, SigningRegion: cfg.AWS.Region}, nil
		})
		options = append(options, awsConfig.WithEndpointResolverWithOptions(resolver))
	}
	if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		options = append(options, awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			cfg.AWS.SessionToken,
		)))
	}
	return awsConfig.LoadDefaultConfig(ctx, options...)
}

func buildAuthProvider(ctx context.Context, cfg *config.Config, awsReady bool) (port.AuthProvider, error) {
	switch cfg.Auth.Mode {
	case "local":
		return auth.NewLocalAuthProvider(cfg.Local.AuthUsers), nil
	default:
		if !awsReady {
			return nil, fmt.Errorf("aws config unavailable for cognito")
		}
		return auth.NewCognitoAuthProvider(ctx, cfg.AWS.Region, cfg.AWS.Endpoint, cfg.AWS.Cognito.ClientID, cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.SessionToken)
	}
}
