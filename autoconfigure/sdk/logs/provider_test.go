package logs_test

import (
	"context"
	otel "github.com/agoda-com/opentelemetry-logs-go"
	autosdk "github.com/agoda-com/opentelemetry-logs-go/autoconfigure/sdk/logs"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"testing"
	"time"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

func doSomething() {
	logger := otel.GetLoggerProvider().Logger(
		instrumentationName,
		logs.WithInstrumentationVersion(instrumentationVersion),
		logs.WithSchemaURL(semconv.SchemaURL),
	)

	body := "My message"
	now := time.Now()
	sn := logs.INFO
	cfg := logs.LogRecordConfig{
		Timestamp:         &now,
		ObservedTimestamp: now,
		Body:              &body,
		SeverityNumber:    &sn,
	}
	logRecord := logs.NewLogRecord(cfg)
	logger.Emit(logRecord)
}
func TestProvider(t *testing.T) {
	ctx := context.Background()
	provider := autosdk.NewLoggerProvider(ctx)
	defer func(provider *sdk.LoggerProvider, ctx context.Context) {
		err := provider.Shutdown(ctx)
		if err != nil {

		}
	}(provider, ctx)

	otel.SetLoggerProvider(provider)

	doSomething()

	assert.NoError(t, provider.Shutdown(ctx))
}
