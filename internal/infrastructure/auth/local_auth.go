package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/config"
	"github.com/vasconcellos/financial-control/internal/domain/port"
)

type LocalAuthProvider struct {
    users       map[string]config.LocalAuthUser
    sessions    map[string]sessionEntry
    sessionTTL  time.Duration
    mu          sync.RWMutex
}

type sessionEntry struct {
    user    config.LocalAuthUser
    expires time.Time
}

var _ port.AuthProvider = (*LocalAuthProvider)(nil)

func NewLocalAuthProvider(users []config.LocalAuthUser) *LocalAuthProvider {
    userMap := make(map[string]config.LocalAuthUser)
    for _, user := range users {
        userMap[user.Username] = user
    }
    return &LocalAuthProvider{
        users:      userMap,
        sessions:   map[string]sessionEntry{},
        sessionTTL: time.Hour,
    }
}

func (p *LocalAuthProvider) Login(ctx context.Context, credentials port.AuthCredentials) (*port.AuthTokens, error) {
	_ = ctx
	user, ok := p.users[credentials.Username]
	if !ok || user.Password != credentials.Password {
		return nil, appErrors.ErrInvalidInput
	}

    accessToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
    refreshToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
    idToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s|%s", user.Email, user.CognitoSub)))

    p.mu.Lock()
    p.cleanupExpiredLocked()
    p.sessions[accessToken] = sessionEntry{user: user, expires: time.Now().Add(p.sessionTTL)}
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
    entry, ok := p.sessions[accessToken]
    p.mu.RUnlock()
    if !ok {
        return nil, fmt.Errorf("invalid token")
    }

    if time.Now().After(entry.expires) {
        p.mu.Lock()
        delete(p.sessions, accessToken)
        p.mu.Unlock()
        return nil, fmt.Errorf("invalid token")
    }

    return map[string]any{
        "username":        entry.user.Username,
        "email":           entry.user.Email,
        "name":            entry.user.Name,
        "defaultCurrency": entry.user.DefaultCurrency,
        "cognitoSub":      entry.user.CognitoSub,
    }, nil
}

func (p *LocalAuthProvider) cleanupExpiredLocked() {
    now := time.Now()
    for token, entry := range p.sessions {
        if now.After(entry.expires) {
            delete(p.sessions, token)
        }
    }
}
