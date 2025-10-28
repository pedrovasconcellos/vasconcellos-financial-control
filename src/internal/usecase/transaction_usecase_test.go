package usecase

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vasconcellos/financial-control/src/internal/domain/dto"
	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
	domainerrors "github.com/vasconcellos/financial-control/src/internal/domain/errors"
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

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, queue, nil, "financial-queue", nil)

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

	// Wait a bit for async goroutine to execute
	time.Sleep(50 * time.Millisecond)

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

func TestTransactionUseCaseAttachReceiptTooLarge(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	txRepo.storage["txn"] = &entity.Transaction{ID: "txn", UserID: "user"}
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{}
	storage := &objectStorageStub{}

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, nil, storage, "queue", nil)

	tooLarge := bytes.Repeat([]byte("a"), int(MaxReceiptSizeBytes)+1)
	_, err := uc.AttachReceipt(context.Background(), "user", "txn", "huge.pdf", "application/pdf", bytes.NewReader(tooLarge))
	if err == nil {
		t.Fatalf("esperava erro para recibo acima do limite")
	}
	if !errors.Is(err, domainerrors.ErrPayloadTooLarge) {
		t.Fatalf("esperava ErrPayloadTooLarge, recebeu %v", err)
	}
	if len(storage.objectKeys) != 0 {
		t.Fatalf("upload não deveria ocorrer quando o recibo excede o limite")
	}
}

func TestTransactionUseCaseRecordTransactionEncryptsNotes(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{categories: map[string]*entity.Category{
		"cat": {ID: "cat", Type: entity.CategoryTypeExpense},
	}}
	queue := &queuePublisherStub{}
	encryptionKey := bytes.Repeat([]byte{1}, 32)

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, queue, nil, "financial-queue", encryptionKey)

	resp, err := uc.RecordTransaction(context.Background(), "user", dto.CreateTransactionRequest{
		AccountID:  "acc",
		CategoryID: "cat",
		Amount:     50,
		Currency:   "USD",
		OccurredAt: time.Now(),
		Notes:      "conteúdo sensível",
	})
	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}
	if resp.Notes != "conteúdo sensível" {
		t.Fatalf("notas deveriam ser retornadas em texto puro")
	}
	stored := txRepo.storage[resp.ID]
	if stored == nil {
		t.Fatalf("transação não armazenada")
	}
	if stored.Notes == "conteúdo sensível" {
		t.Fatalf("notas deveriam estar criptografadas")
	}
	if stored.Metadata["notes_encrypted"] != "true" {
		t.Fatalf("flag de notas criptografadas deveria ser verdadeira")
	}

	txRepo.listResponse = []*entity.Transaction{stored}
	list, err := uc.ListTransactions(context.Background(), "user", time.Now().Add(-time.Hour), time.Now().Add(time.Hour), 10, 0)
	if err != nil {
		t.Fatalf("não esperava erro ao listar: %v", err)
	}
	if len(list) != 1 || list[0].Notes != "conteúdo sensível" {
		t.Fatalf("notas deveriam ser decriptografadas na resposta")
	}
}

func TestTransactionUseCaseUpdateTransactionEncryptsNotes(t *testing.T) {
	txRepo := newTransactionRepositoryStub()
	transaction := &entity.Transaction{
		ID:         "txn",
		UserID:     "user",
		AccountID:  "acc",
		CategoryID: "cat",
		Notes:      "plain",
		Metadata:   map[string]string{},
	}
	txRepo.storage["txn"] = transaction
	accountRepo := newAccountRepositoryStub()
	categoryRepo := &categoryRepositoryStub{}
	encryptionKey := bytes.Repeat([]byte{2}, 32)

	uc := NewTransactionUseCase(txRepo, accountRepo, categoryRepo, nil, nil, "queue", encryptionKey)

	newNotes := "nota atualizada"
	resp, err := uc.UpdateTransaction(context.Background(), "user", "txn", dto.UpdateTransactionRequest{
		Notes: &newNotes,
	})
	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}
	if resp.Notes != newNotes {
		t.Fatalf("notas deveriam ser retornadas em texto puro após atualização")
	}
	if txRepo.lastUpdated == nil {
		t.Fatalf("update deveria ser chamado")
	}
	if txRepo.lastUpdated.Notes == newNotes {
		t.Fatalf("notas persistidas deveriam estar criptografadas")
	}
	if txRepo.lastUpdated.Metadata["notes_encrypted"] != "true" {
		t.Fatalf("flag de notas criptografadas deveria ser atualizada")
	}
}
