package usecase

import (
	"context"
	"io"
	"time"

	"github.com/vasconcellos/financial-control/internal/domain/entity"
	"github.com/vasconcellos/financial-control/internal/domain/port"
)

type accountRepositoryStub struct {
	created     []*entity.Account
	storage     map[string]*entity.Account
	adjustments []float64
	lastLimit   int64
	lastOffset  int64
}

func newAccountRepositoryStub() *accountRepositoryStub {
	return &accountRepositoryStub{storage: make(map[string]*entity.Account)}
}

func (s *accountRepositoryStub) Create(ctx context.Context, account *entity.Account) error {
	s.created = append(s.created, account)
	s.storage[account.ID] = account
	return nil
}

func (s *accountRepositoryStub) Update(ctx context.Context, account *entity.Account) error {
	s.storage[account.ID] = account
	return nil
}

func (s *accountRepositoryStub) Delete(ctx context.Context, id string, userID string) error {
	delete(s.storage, id)
	return nil
}

func (s *accountRepositoryStub) GetByID(ctx context.Context, id string, userID string) (*entity.Account, error) {
	return s.storage[id], nil
}

func (s *accountRepositoryStub) List(ctx context.Context, userID string, limit int64, offset int64) ([]*entity.Account, error) {
	s.lastLimit = limit
	s.lastOffset = offset
	var filtered []*entity.Account
	for _, account := range s.storage {
		if account.UserID == userID {
			filtered = append(filtered, account)
		}
	}

	start := offset
	if start > int64(len(filtered)) {
		start = int64(len(filtered))
	}
	end := start + limit
	if limit <= 0 || end > int64(len(filtered)) {
		end = int64(len(filtered))
	}

	startIdx := int(start)
	endIdx := int(end)
	return filtered[startIdx:endIdx], nil
}

func (s *accountRepositoryStub) AdjustBalance(ctx context.Context, id string, userID string, amount float64) error {
	s.adjustments = append(s.adjustments, amount)
	if acc, ok := s.storage[id]; ok {
		acc.Balance += amount
	}
	return nil
}

type transactionRepositoryStub struct {
	created      []*entity.Transaction
	storage      map[string]*entity.Transaction
	lastUpdated  *entity.Transaction
	listResponse []*entity.Transaction
	lastLimit    int64
	lastOffset   int64
}

func newTransactionRepositoryStub() *transactionRepositoryStub {
	return &transactionRepositoryStub{storage: make(map[string]*entity.Transaction)}
}

func (s *transactionRepositoryStub) Create(ctx context.Context, transaction *entity.Transaction) error {
	s.created = append(s.created, transaction)
	s.storage[transaction.ID] = transaction
	return nil
}

func (s *transactionRepositoryStub) Update(ctx context.Context, transaction *entity.Transaction) error {
	s.storage[transaction.ID] = transaction
	s.lastUpdated = transaction
	return nil
}

func (s *transactionRepositoryStub) GetByID(ctx context.Context, id string, userID string) (*entity.Transaction, error) {
	return s.storage[id], nil
}

func (s *transactionRepositoryStub) List(ctx context.Context, userID string, from time.Time, to time.Time, limit int64, offset int64) ([]*entity.Transaction, error) {
	if s.listResponse != nil {
		return s.listResponse, nil
	}
	var result []*entity.Transaction
	for _, txn := range s.storage {
		if txn.UserID == userID {
			result = append(result, txn)
		}
	}
	s.lastLimit = limit
	s.lastOffset = offset
	start := offset
	if start > int64(len(result)) {
		start = int64(len(result))
	}
	end := start + limit
	if limit <= 0 || end > int64(len(result)) {
		end = int64(len(result))
	}
	startIdx := int(start)
	endIdx := int(end)
	return result[startIdx:endIdx], nil
}

func (s *transactionRepositoryStub) ListByCategory(ctx context.Context, userID string, categoryID string, from time.Time, to time.Time) ([]*entity.Transaction, error) {
	return nil, nil
}

type categoryRepositoryStub struct {
	categories map[string]*entity.Category
}

func (s *categoryRepositoryStub) Create(ctx context.Context, category *entity.Category) error {
	return nil
}

func (s *categoryRepositoryStub) Update(ctx context.Context, category *entity.Category) error {
	return nil
}

func (s *categoryRepositoryStub) Delete(ctx context.Context, id string, userID string) error {
	return nil
}

func (s *categoryRepositoryStub) GetByID(ctx context.Context, id string, userID string) (*entity.Category, error) {
	if s.categories == nil {
		return nil, nil
	}
	return s.categories[id], nil
}

func (s *categoryRepositoryStub) List(ctx context.Context, userID string) ([]*entity.Category, error) {
	return nil, nil
}

type queuePublisherStub struct {
	lastMessage port.QueueMessage
	called      bool
}

func (s *queuePublisherStub) Publish(ctx context.Context, queueName string, message port.QueueMessage) error {
	s.called = true
	s.lastMessage = message
	return nil
}

type objectStorageStub struct {
	objectKeys []string
}

func (s *objectStorageStub) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	s.objectKeys = append(s.objectKeys, key)
	return key, nil
}

func (s *objectStorageStub) GetPresignedURL(ctx context.Context, key string) (string, error) {
	return "https://example.com/" + key, nil
}
