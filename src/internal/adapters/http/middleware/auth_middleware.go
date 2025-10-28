package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vasconcellos/financial-control/src/internal/usecase"
)

type AuthMiddleware struct {
	authUseCase *usecase.AuthUseCase
	userUseCase *usecase.UserUseCase
}

func NewAuthMiddleware(authUseCase *usecase.AuthUseCase, userUseCase *usecase.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{authUseCase: authUseCase, userUseCase: userUseCase}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := m.authUseCase.ValidateToken(c.Request.Context(), tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		currency, _ := claims["custom:defaultCurrency"].(string)
		if currency == "" {
			currency, _ = claims["defaultCurrency"].(string)
		}
		if currency == "" {
			currency = "USD"
		}
		sub, _ := claims["sub"].(string)
		if sub == "" {
			sub, _ = claims["cognito:username"].(string)
		}
		if sub == "" {
			sub, _ = claims["cognitoSub"].(string)
		}

		user, err := m.userUseCase.EnsureUser(c.Request.Context(), email, name, sub, currency)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unable to load user"})
			return
		}

		SetUserContext(c, AuthenticatedUser{
			ID:              user.ID,
			Email:           user.Email,
			Name:            user.Name,
			DefaultCurrency: user.DefaultCurrency.String(),
			CognitoSub:      user.CognitoSub,
		})

		c.Next()
	}
}
