package logs

import (
	"github.com/agoda-com/opentelemetry-logs-go/otel/logs"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"testing"
	"time"
)

func TestLogsReadWriteAPIFormat(t *testing.T) {

	traceID, _ := trace.TraceIDFromHex("80f198ee56343ba864fe8b2a57d3eff7")
	spanID, _ := trace.SpanIDFromHex("2a00000000000000")

	spanCtxCfg := trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
		Remote:     false,
	}

	spanCtx := trace.NewSpanContext(spanCtxCfg)

	body := "My Log Message"
	severityText := "INFO"
	severityNumber := logs.INFO
	resource := sdkresource.NewWithAttributes("http", semconv.HTTPURL("testurl"))
	instrumentationScope := instrumentation.Scope{
		Name:      "test",
		Version:   "testVersion",
		SchemaURL: "http",
	}
	attributes := []attribute.KeyValue{semconv.HTTPURL("testurl")}
	timestamp := time.Now()

	record := newReadWriteLogRecord(
		&spanCtx,
		&body,
		&severityText,
		&severityNumber,
		resource,
		&instrumentationScope,
		&attributes,
		&timestamp,
	)

	assert.Equal(t, "My Log Message", *record.Body())

}
