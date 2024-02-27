package stdoutlogs

import (
	"bytes"
	"context"
	"io"
	"log"
	"testing"
	"time"

	otel "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otlplogs-example"),
		semconv.ServiceVersion("0.0.1"),
	)
}

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

func installExportPipeline(writer io.Writer) (func(context.Context) error, error) {
	exporter, _ := NewExporter(WithWriter(writer))

	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithSyncer(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}

func TestStdoutExporter(t *testing.T) {
	{
		ctx := context.Background()

		var writer bytes.Buffer
		// Registers a logger Provider globally.
		shutdown, err := installExportPipeline(&writer)
		if err != nil {
			log.Fatal(err)
		}
		doSomething()

		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		actual := writer.String()

		assert.Contains(t, actual, "My message")
		assert.Contains(t, actual, "otlplogs-example")
		assert.Contains(t, actual, "0.0.1")
	}
}
