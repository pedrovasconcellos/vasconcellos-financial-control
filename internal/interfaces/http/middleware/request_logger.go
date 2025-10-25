package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware adiciona identificador e métricas de requisição
type RequestLoggerMiddleware struct {
	logger *zap.Logger
}

func NewRequestLoggerMiddleware(logger *zap.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{logger: logger}
}

func (m *RequestLoggerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		start := time.Now()

		requestLogger := m.logger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("client_ip", c.ClientIP()),
		)

		attachLogger(c, requestLogger)
		requestLogger.Info("request started")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		requestLogger.Info("request completed",
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
		)
	}
}
