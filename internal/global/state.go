package global

import (
	"errors"
	"github.com/agoda-com/otellogs-go/otel/logs"
	"sync"
	"sync/atomic"
)

type (
	loggerProviderHolder struct {
		lp logs.LoggerProvider
	}
)

var (
	globalOtelLogger = defaultLoggerValue()

	delegateLoggerOnce sync.Once
)

// LoggerProvider is the internal implementation for global.LoggerProvider.
func LoggerProvider() logs.LoggerProvider {
	return globalOtelLogger.Load().(loggerProviderHolder).lp
}

// SetLoggerProvider is the internal implementation for global.SetLoggerProvider.
func SetLoggerProvider(lp logs.LoggerProvider) {
	current := LoggerProvider()

	if _, cOk := current.(*loggerProvider); cOk {
		if _, tpOk := lp.(*loggerProvider); tpOk && current == lp {
			// Do not assign the default delegating LoggerProvider to delegate
			// to itself.
			Error(
				errors.New("no delegate configured in logger provider"),
				"Setting logger provider to it's current value. No delegate will be configured",
			)
			return
		}
	}

	globalOtelLogger.Store(loggerProviderHolder{lp: lp})
}

func defaultLoggerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(loggerProviderHolder{lp: &loggerProvider{}})
	return v
}
