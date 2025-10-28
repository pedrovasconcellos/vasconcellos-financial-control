package middleware

import "github.com/gin-gonic/gin"

type AuthenticatedUser struct {
	ID              string
	Email           string
	Name            string
	DefaultCurrency string
	CognitoSub      string
}

const userContextKey = "authenticatedUser"

func SetUserContext(c *gin.Context, user AuthenticatedUser) {
	c.Set(userContextKey, user)
}

func GetUserContext(c *gin.Context) (AuthenticatedUser, bool) {
	value, exists := c.Get(userContextKey)
	if !exists {
		return AuthenticatedUser{}, false
	}
	user, ok := value.(AuthenticatedUser)
	return user, ok
}
