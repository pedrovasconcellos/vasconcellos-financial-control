package usecase

import (
	"context"
	"testing"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/entity"
)

// TestAccountUseCaseCreate lista a criação básica garantindo persistência no repositório
func TestAccountUseCaseCreate(t *testing.T) {
	repo := newAccountRepositoryStub()
	uc := NewAccountUseCase(repo)

	resp, err := uc.CreateAccount(context.Background(), "user-1", dto.CreateAccountRequest{
		Name:     "Main",
		Type:     "checking",
		Currency: "USD",
	})
	if err != nil {
		t.Fatalf("esperava criação sem erros, obteve: %v", err)
	}
	if resp.ID == "" {
		t.Fatalf("ID não deveria estar vazio")
	}
	if len(repo.created) != 1 {
		t.Fatalf("esperava 1 conta criada, obtido %d", len(repo.created))
	}
	if repo.created[0].Name != "Main" {
		t.Errorf("nome inesperado: %s", repo.created[0].Name)
	}
}

// TestAccountUseCaseList garante que o caso de uso devolve todas as contas do usuário
func TestAccountUseCaseList(t *testing.T) {
	repo := newAccountRepositoryStub()
	repo.Create(context.Background(), &entity.Account{ID: "a1", UserID: "user-1", Name: "Conta A"})
	repo.Create(context.Background(), &entity.Account{ID: "a2", UserID: "user-1", Name: "Conta B"})
	repo.Create(context.Background(), &entity.Account{ID: "a3", UserID: "user-2", Name: "Conta C"})

	uc := NewAccountUseCase(repo)
	resp, err := uc.ListAccounts(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("esperava listagem sem erros, obteve: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("esperava 2 contas, obtido %d", len(resp))
	}
}
