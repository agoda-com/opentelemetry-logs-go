package otel

import (
	"github.com/agoda-com/opentelemetry-logs-go/otel/internal/global"
	"github.com/agoda-com/opentelemetry-logs-go/otel/logs"
)

// GetLoggerProvider returns the registered global logger provider.
// If none is registered then an instance of NoopLoggerProvider is returned.
//
// loggerProvider := otel.GetLoggerProvider()
func GetLoggerProvider() logs.LoggerProvider {
	return global.LoggerProvider()
}

// SetLoggerProvider registers `lp` as the global logger provider.
func SetLoggerProvider(lp logs.LoggerProvider) {
	global.SetLoggerProvider(lp)
}
