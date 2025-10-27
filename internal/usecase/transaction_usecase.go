package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/financial-control/internal/domain/dto"
	"github.com/vasconcellos/financial-control/internal/domain/entity"
	"github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/domain/port"
	"github.com/vasconcellos/financial-control/internal/domain/repository"
	"github.com/vasconcellos/financial-control/internal/infrastructure/security"
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
	queuePublisher  port.QueuePublisher
	storage         port.ObjectStorage
	eventQueueName  string
	encryptionKey   []byte
}

const MaxReceiptSizeBytes int64 = 5 * 1024 * 1024
const notesEncryptedMetadataKey = "notes_encrypted"

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	queuePublisher port.QueuePublisher,
	storage port.ObjectStorage,
	eventQueueName string,
	encryptionKey []byte,
) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		queuePublisher:  queuePublisher,
		storage:         storage,
		eventQueueName:  eventQueueName,
		encryptionKey:   encryptionKey,
	}
}

func (uc *TransactionUseCase) encryptNotes(notes string, metadata map[string]string) (string, error) {
	if notes == "" {
		if metadata != nil {
			delete(metadata, notesEncryptedMetadataKey)
		}
		return "", nil
	}
	if len(uc.encryptionKey) == 0 {
		if metadata != nil {
			delete(metadata, notesEncryptedMetadataKey)
		}
		return notes, nil
	}
	if metadata != nil {
		metadata[notesEncryptedMetadataKey] = "true"
	}

	ciphertext, err := security.EncryptAESGCM([]byte(notes), uc.encryptionKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (uc *TransactionUseCase) decryptNotes(notes string, metadata map[string]string) (string, error) {
	if notes == "" {
		return "", nil
	}
	if len(uc.encryptionKey) == 0 {
		return notes, nil
	}
	if metadata == nil || metadata[notesEncryptedMetadataKey] != "true" {
		return notes, nil
	}

	raw, err := base64.StdEncoding.DecodeString(notes)
	if err != nil {
		return "", err
	}
	plaintext, err := security.DecryptAESGCM(raw, uc.encryptionKey)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
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
		Currency:    entity.Currency(request.Currency),
		Description: request.Description,
		OccurredAt:  request.OccurredAt,
		Status:      entity.TransactionStatusCompleted,
		Tags:        request.Tags,
		Notes:       request.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    map[string]string{},
	}
	encryptedNotes, err := uc.encryptNotes(transaction.Notes, transaction.Metadata)
	if err != nil {
		return nil, err
	}
	transaction.Notes = encryptedNotes

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
	notesValue, err := uc.decryptNotes(transaction.Notes, transaction.Metadata)
	if err != nil {
		return nil, err
	}

	if uc.queuePublisher != nil && uc.eventQueueName != "" {
		// Publish to SQS asynchronously to avoid blocking HTTP response
		go func() {
			messagePayload := map[string]any{
				"transactionId": transaction.ID,
				"userId":        transaction.UserID,
				"occurredAt":    transaction.OccurredAt,
				"amount":        transaction.Amount,
				"currency":      transaction.Currency.String(),
				"categoryId":    transaction.CategoryID,
				"accountId":     transaction.AccountID,
				"type":          category.Type,
			}
			body, marshalErr := json.Marshal(messagePayload)
			if marshalErr == nil {
				_ = uc.queuePublisher.Publish(context.Background(), uc.eventQueueName, port.QueueMessage{
					ID:         transaction.ID,
					Payload:    body,
					Attributes: map[string]string{"eventType": "TRANSACTION_RECORDED"},
				})
			}
		}()
	}

	return &dto.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency.String(),
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       notesValue,
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
	if request.Notes != nil {
		if transaction.Metadata == nil {
			transaction.Metadata = map[string]string{}
		}
		encryptedNotes, encryptErr := uc.encryptNotes(transaction.Notes, transaction.Metadata)
		if encryptErr != nil {
			return nil, encryptErr
		}
		transaction.Notes = encryptedNotes
	}

	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}
	notesValue, err := uc.decryptNotes(transaction.Notes, transaction.Metadata)
	if err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency.String(),
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       notesValue,
	}, nil
}

func (uc *TransactionUseCase) ListTransactions(ctx context.Context, userID string, from, to time.Time, limit int64, offset int64) ([]*dto.TransactionResponse, error) {
	transactions, err := uc.transactionRepo.List(ctx, userID, from, to, limit, offset)
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

		notesValue, notesErr := uc.decryptNotes(transaction.Notes, transaction.Metadata)
		if notesErr != nil {
			return nil, notesErr
		}

		response = append(response, &dto.TransactionResponse{
			ID:          transaction.ID,
			AccountID:   transaction.AccountID,
			CategoryID:  transaction.CategoryID,
			Amount:      transaction.Amount,
			Currency:    transaction.Currency.String(),
			Description: transaction.Description,
			OccurredAt:  transaction.OccurredAt,
			Status:      string(transaction.Status),
			Tags:        transaction.Tags,
			Notes:       notesValue,
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

	tempFile, err := os.CreateTemp("", "receipt-*")
	if err != nil {
		return nil, err
	}
	defer func() {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	limited := &io.LimitedReader{R: data, N: MaxReceiptSizeBytes + 1}
	written, err := io.Copy(tempFile, limited)
	if err != nil {
		return nil, err
	}
	if written == 0 {
		return nil, fmt.Errorf("empty receipt")
	}
	if written > MaxReceiptSizeBytes {
		return nil, errors.ErrPayloadTooLarge
	}

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	objectKey := fmt.Sprintf("users/%s/transactions/%s/%s", userID, transactionID, filename)
	if _, err := uc.storage.Upload(ctx, objectKey, tempFile, contentType); err != nil {
		return nil, err
	}
	transaction.ReceiptObject = &objectKey
	transaction.UpdatedAt = time.Now().UTC()

	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}
	notesValue, err := uc.decryptNotes(transaction.Notes, transaction.Metadata)
	if err != nil {
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
		Currency:    transaction.Currency.String(),
		Description: transaction.Description,
		OccurredAt:  transaction.OccurredAt,
		Status:      string(transaction.Status),
		Tags:        transaction.Tags,
		Notes:       notesValue,
	}
	response.ReceiptURL = &url

	return response, nil
}
