package logs

// NewNoopLoggerProvider returns an implementation of LoggerProvider that
// performs no operations. The Logger created from the returned
// LoggerProvider also perform no operations.
func NewNoopLoggerProvider() LoggerProvider {
	return noopLoggerProvider{}
}

type noopLoggerProvider struct{}

var _ LoggerProvider = noopLoggerProvider{}

func (p noopLoggerProvider) Logger(string, ...LoggerOption) Logger {
	return noopLogger{}
}

type noopLogger struct{}

var _ Logger = noopLogger{}

func (n noopLogger) Emit(logRecord LogRecord) {}
