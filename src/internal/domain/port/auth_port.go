package port

import "context"

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int32
	TokenType    string
}

type AuthCredentials struct {
	Username string
	Password string
}

type AuthProvider interface {
	Login(ctx context.Context, credentials AuthCredentials) (*AuthTokens, error)
	Validate(ctx context.Context, accessToken string) (map[string]any, error)
}
