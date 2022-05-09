package utils

import (
	"go.uber.org/zap"
)

var zapLogger *zap.Logger
func NewZapLogger(app string) *zap.Logger {
	if zapLogger != nil {
		return zapLogger
	}
	if IsDev() {
		zapLogger, _ = zap.NewDevelopment()
	} else {
		zapLogger, _ = zap.NewProduction()
	}
	zapLogger = zapLogger.With(zap.String("app", app))
	return zapLogger
}