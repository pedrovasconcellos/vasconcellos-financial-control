package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const loggerContextKey = "request_logger"

// recupera o logger do contexto atual ou devolve um no-op
func LoggerFromContext(c *gin.Context) *zap.Logger {
	if value, exists := c.Get(loggerContextKey); exists {
		if logger, ok := value.(*zap.Logger); ok && logger != nil {
			return logger
		}
	}
	return zap.NewNop()
}

func attachLogger(c *gin.Context, logger *zap.Logger) {
	c.Set(loggerContextKey, logger)
}
