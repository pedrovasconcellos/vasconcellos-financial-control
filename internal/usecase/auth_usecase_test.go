package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/domain/port"
)

type fakeAuthProvider struct {
	tokens *port.AuthTokens
	err    error
}

func (f *fakeAuthProvider) Login(ctx context.Context, credentials port.AuthCredentials) (*port.AuthTokens, error) {
	return f.tokens, f.err
}

func (f *fakeAuthProvider) Validate(ctx context.Context, token string) (map[string]any, error) {
	return nil, errors.New("not implemented")
}

// TestAuthUseCaseLoginSucesso garante que o caso de uso propaga os tokens recebidos do provedor
func TestAuthUseCaseLoginSucesso(t *testing.T) {
	expected := &port.AuthTokens{
		AccessToken:  "access",
		RefreshToken: "refresh",
		IDToken:      "id",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}
	uc := NewAuthUseCase(&fakeAuthProvider{tokens: expected})

	resp, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "user@test.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}
	if resp.AccessToken != expected.AccessToken {
		t.Errorf("token inesperado: %s", resp.AccessToken)
	}
	if resp.ExpiresIn != expected.ExpiresIn {
		t.Errorf("expiresIn inesperado: %d", resp.ExpiresIn)
	}
}

// TestAuthUseCaseLoginErro garante que um erro do provedor é repassado ao chamador
func TestAuthUseCaseLoginErro(t *testing.T) {
	uc := NewAuthUseCase(&fakeAuthProvider{err: errors.New("invalid credentials")})

	_, err := uc.Login(context.Background(), dto.LoginRequest{})
	if err == nil {
		t.Fatalf("esperava erro, obteve nil")
	}
}
