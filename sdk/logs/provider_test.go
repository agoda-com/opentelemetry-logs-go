package logs

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"testing"
)

const (
	instrumentationName    = "otel/zap"
	instrumentationVersion = "0.0.1"
)

func TestLogsProvider(t *testing.T) {

	noopProvider := logs.NewNoopLoggerProvider()

	noopLogger := noopProvider.Logger(instrumentationName,
		logs.WithInstrumentationVersion(instrumentationVersion),
		logs.WithSchemaURL(semconv.SchemaURL),
	)

	logsOtlpExporter, _ := otlplogshttp.New(context.Background())

	batchOtlpLogger := NewLoggerProvider(
		WithLogRecordProcessor(NewBatchLogRecordProcessor(logsOtlpExporter)),
		WithResource(
			resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("unit_test"))),
	).Logger("otlp",
		logs.WithInstrumentationVersion(instrumentationName),
		logs.WithSchemaURL(instrumentationVersion),
		logs.WithInstrumentationAttributes(semconv.HostName("some.host")),
	)

	body := "body"

	logRecord := logs.NewLogRecord(logs.LogRecordConfig{
		Body: &body,
	})

	noopLogger.Emit(logRecord)
	batchOtlpLogger.Emit(logRecord)

}
