package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/entity"
	"github.com/vasconcellos/finance-control/internal/domain/errors"
	"github.com/vasconcellos/finance-control/internal/domain/port"
	"github.com/vasconcellos/finance-control/internal/domain/repository"
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
	queuePublisher  port.QueuePublisher
	storage         port.ObjectStorage
	eventQueueName  string
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	queuePublisher port.QueuePublisher,
	storage port.ObjectStorage,
	eventQueueName string,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		queuePublisher:  queuePublisher,
		storage:         storage,
		eventQueueName:  eventQueueName,
	}
}

func (uc *TransactionUseCase) RecordTransaction(ctx context.Context, userID string, request dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	if request.Amount <= 0 {
		return nil, errors.ErrInvalidInput
	}
	category, err := uc.categoryRepo.GetByID(ctx, request.CategoryID, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.ErrInvalidInput
	}

	now := time.Now().UTC()
	transaction := &entity.Transaction{
		ID:          uuid.NewString(),
		UserID:      userID,
		AccountID:   request.AccountID,
		CategoryID:  request.CategoryID,
		Amount:      request.Amount,
		Currency:    request.Currency,
		Description: request.Description,
		OccurredAt:  request.OccurredAt,
		Status:      entity.TransactionStatusCompleted,
		Tags:        request.Tags,
		Notes:       request.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    map[string]string{},
	}

	if category.Type == entity.CategoryTypeExpense {
		if err := uc.accountRepo.AdjustBalance(ctx, request.AccountID, userID, -request.Amount); err != nil {
			return nil, err
		}
	} else {
		if err := uc.accountRepo.AdjustBalance(ctx, request.AccountID, userID, request.Amount); err != nil {
			return nil, err
		}
	}

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	if uc.queuePublisher != nil && uc.eventQueueName != "" {
		messagePayload := map[string]any{
			"transactionId": transaction.ID,
			"userId":        transaction.UserID,
			"occurredAt":    transaction.OccurredAt,
			"amount":        transaction.Amount,
			"currency":      transaction.Currency,
			"categoryId":    transaction.CategoryID,
			"accountId":     transaction.AccountID,
			"type":          category.Type,
		}
		body, marshalErr := json.Marshal(messagePayload)
		if marshalErr == nil {
			_ = uc.queuePublisher.Publish(ctx, uc.eventQueueName, port.QueueMessage{
				ID:         transaction.ID,
				Payload:    body,
				Attributes: map[string]string{"eventType": "TRANSACTION_RECORDED"},
			})
		}
	}

	return &dto.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       transaction.Notes,
	}, nil
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, userID string, transactionID string, request dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID, userID)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.ErrNotFound
	}

	if request.CategoryID != nil {
		transaction.CategoryID = *request.CategoryID
	}
	if request.Description != nil {
		transaction.Description = *request.Description
	}
	if request.Notes != nil {
		transaction.Notes = *request.Notes
	}
	if request.Tags != nil {
		transaction.Tags = request.Tags
	}
	if request.Status != nil {
		transaction.Status = entity.TransactionStatus(*request.Status)
	}
	transaction.UpdatedAt = time.Now().UTC()

	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       transaction.Notes,
	}, nil
}

func (uc *TransactionUseCase) ListTransactions(ctx context.Context, userID string, from, to time.Time) ([]*dto.TransactionResponse, error) {
	transactions, err := uc.transactionRepo.List(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.TransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		var receiptURL *string
		if transaction.ReceiptObject != nil && uc.storage != nil {
			url, getErr := uc.storage.GetPresignedURL(ctx, *transaction.ReceiptObject)
			if getErr == nil {
				receiptURL = &url
			}
		}

		response = append(response, &dto.TransactionResponse{
			ID:          transaction.ID,
			AccountID:   transaction.AccountID,
			CategoryID:  transaction.CategoryID,
			Amount:      transaction.Amount,
			Currency:    transaction.Currency,
			Description: transaction.Description,
			OccurredAt:  transaction.OccurredAt,
			Status:      string(transaction.Status),
			Tags:        transaction.Tags,
			Notes:       transaction.Notes,
			ReceiptURL:  receiptURL,
		})
	}

	return response, nil
}

func (uc *TransactionUseCase) AttachReceipt(ctx context.Context, userID string, transactionID string, filename string, contentType string, data io.Reader) (*dto.TransactionResponse, error) {
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID, userID)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.ErrNotFound
	}

	if uc.storage == nil {
		return nil, fmt.Errorf("object storage disabled")
	}

	objectKey := fmt.Sprintf("users/%s/transactions/%s/%s", userID, transactionID, filename)
	if _, err := uc.storage.Upload(ctx, objectKey, data, contentType); err != nil {
		return nil, err
	}
	transaction.ReceiptObject = &objectKey
	transaction.UpdatedAt = time.Now().UTC()

	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	url, err := uc.storage.GetPresignedURL(ctx, objectKey)
	if err != nil {
		return nil, err
	}

	response := &dto.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       transaction.Notes,
	}
	response.ReceiptURL = &url

	return response, nil
}
