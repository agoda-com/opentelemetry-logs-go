package otel

import (
	"github.com/agoda-com/opentelemetry-logs-go/otel/logs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testLoggerProvider struct{}

var _ logs.LoggerProvider = &testLoggerProvider{}

func (*testLoggerProvider) Logger(_ string, _ ...logs.LoggerOption) logs.Logger {
	return logs.NewNoopLoggerProvider().Logger("")
}

func TestMultipleGlobalTracerProvider(t *testing.T) {
	p1 := testLoggerProvider{}
	p2 := logs.NewNoopLoggerProvider()
	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	got := GetLoggerProvider()
	assert.Equal(t, p2, got)
}
