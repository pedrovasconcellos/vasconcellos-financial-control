package usecase

import (
	"context"

	"github.com/vasconcellos/financial-control/src/internal/domain/dto"
	"github.com/vasconcellos/financial-control/src/internal/domain/port"
)

type AuthUseCase struct {
	authProvider port.AuthProvider
}

func NewAuthUseCase(authProvider port.AuthProvider) *AuthUseCase {
	return &AuthUseCase{
		authProvider: authProvider,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	tokens, err := uc.authProvider.Login(ctx, port.AuthCredentials{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IDToken:      tokens.IDToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (map[string]any, error) {
	return uc.authProvider.Validate(ctx, token)
}
