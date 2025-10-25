package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vasconcellos/finance-control/internal/config"
	"github.com/vasconcellos/finance-control/internal/domain/port"
)

type LocalAuthProvider struct {
	users    map[string]config.LocalAuthUser
	sessions map[string]config.LocalAuthUser
	mu       sync.RWMutex
}

var _ port.AuthProvider = (*LocalAuthProvider)(nil)

func NewLocalAuthProvider(users []config.LocalAuthUser) *LocalAuthProvider {
	userMap := make(map[string]config.LocalAuthUser)
	for _, user := range users {
		userMap[user.Username] = user
	}
	return &LocalAuthProvider{
		users:    userMap,
		sessions: map[string]config.LocalAuthUser{},
	}
}

func (p *LocalAuthProvider) Login(ctx context.Context, credentials port.AuthCredentials) (*port.AuthTokens, error) {
	_ = ctx
	user, ok := p.users[credentials.Username]
	if !ok || user.Password != credentials.Password {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
	refreshToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
	idToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s|%s", user.Email, user.CognitoSub)))

	p.mu.Lock()
	p.sessions[accessToken] = user
	p.mu.Unlock()

	return &port.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
		ExpiresIn:    int32(time.Hour.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (p *LocalAuthProvider) Validate(ctx context.Context, accessToken string) (map[string]any, error) {
	_ = ctx
	p.mu.RLock()
	user, ok := p.sessions[accessToken]
	p.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return map[string]any{
		"username":        user.Username,
		"email":           user.Email,
		"name":            user.Name,
		"defaultCurrency": user.DefaultCurrency,
		"cognitoSub":      user.CognitoSub,
	}, nil
}
