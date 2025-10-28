package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(environment string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if environment == "development" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	return cfg.Build()
}
