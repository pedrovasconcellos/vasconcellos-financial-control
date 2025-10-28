package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPSRedirectMiddleware força redirecionamento para HTTPS
type HTTPSRedirectMiddleware struct {
	enabled     bool
	environment string
}

// NewHTTPSRedirectMiddleware cria um novo middleware de redirecionamento HTTPS
func NewHTTPSRedirectMiddleware(enabled bool, environment string) *HTTPSRedirectMiddleware {
	return &HTTPSRedirectMiddleware{
		enabled:     enabled,
		environment: environment,
	}
}

// Handle implementa o middleware de redirecionamento HTTPS
func (m *HTTPSRedirectMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Se HTTPS não está habilitado, não faz nada
		if !m.enabled {
			c.Next()
			return
		}

		// Em desenvolvimento, permite HTTP sem redirecionamento
		if m.environment == "development" {
			c.Next()
			return
		}

		// Verifica se a requisição já é HTTPS
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Next()
			return
		}

		// Para ambientes não-desenvolvimento, força HTTPS
		httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
		c.Redirect(http.StatusMovedPermanently, httpsURL)
		c.Abort()
	}
}
