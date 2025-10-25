package usecase

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/entity"
	domainerrors "github.com/vasconcellos/finance-control/internal/domain/errors"
)

// TestTransactionUseCaseRecordTransactionValorInvalido garante que valores não positivos geram erro de validação
func TestTransactionUseCaseRecordTransactionValorInvalido(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{categories: map[string]*entity.Category{
		"cat": {ID: "cat", Type: entity.CategoryTypeExpense},
	}}

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, nil, nil, "queue", nil)

	_, err := uc.RecordTransaction(context.Background(), "user", dto.CreateTransactionRequest{
		AccountID:  "acc",
		CategoryID: "cat",
		Amount:     0,
		Currency:   "USD",
	})
	if !errors.Is(err, domainerrors.ErrInvalidInput) {
		t.Fatalf("esperava ErrInvalidInput, obteve %v", err)
	}
}

// TestTransactionUseCaseRecordTransactionDespesa garante ajuste negativo no saldo e envio ao SQS
func TestTransactionUseCaseRecordTransactionDespesa(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{categories: map[string]*entity.Category{
		"cat": {ID: "cat", Type: entity.CategoryTypeExpense},
	}}
	queue := &queuePublisherStub{}

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, queue, nil, "finance-queue", nil)

	_, err := uc.RecordTransaction(context.Background(), "user", dto.CreateTransactionRequest{
		AccountID:  "acc",
		CategoryID: "cat",
		Amount:     100,
		Currency:   "USD",
		OccurredAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}
	if len(accountRepo.adjustments) != 1 || accountRepo.adjustments[0] != -100 {
		t.Fatalf("ajuste de saldo inadequado: %#v", accountRepo.adjustments)
	}
	if !queue.called {
		t.Fatalf("publicação no SQS deveria ocorrer")
	}
	if txRepo.created == nil || len(txRepo.created) != 1 {
		t.Fatalf("transação não persistida corretamente")
	}
}

// TestTransactionUseCaseAttachReceipt garante upload e retorno da URL pré-assinada
func TestTransactionUseCaseAttachReceipt(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	transaction := &entity.Transaction{ID: "txn", UserID: "user", ReceiptObject: nil}
	txRepo.storage["txn"] = transaction
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{}
	queue := &queuePublisherStub{}
	storage := &objectStorageStub{}

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, queue, storage, "queue", nil)

	resp, err := uc.AttachReceipt(context.Background(), "user", "txn", "receipt.pdf", "application/pdf", bytes.NewReader([]byte("filedata")))
	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}
	if resp.ReceiptURL == nil || *resp.ReceiptURL == "" {
		t.Fatalf("URL de recibo deveria ser preenchida")
	}
	if len(storage.objectKeys) != 1 {
		t.Fatalf("upload não registrado")
	}
}
