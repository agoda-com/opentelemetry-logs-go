package logstest

import (
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	logssdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"time"
)

// LogRecordStubs is a slice of LogRecordStub use for testing an SDK.
type LogRecordStubs []LogRecordStub

// Snapshots returns s as a slice of ReadOnlySpans.
func (l LogRecordStubs) Snapshots() []logssdk.ReadableLogRecord {
	if len(l) == 0 {
		return nil
	}

	rlr := make([]logssdk.ReadableLogRecord, len(l))
	for i := 0; i < len(l); i++ {
		rlr[i] = l[i].Snapshot()
	}
	return rlr
}

// LogRecordStub is a stand-in for a LogRecord.
type LogRecordStub struct {
	Timestamp            *time.Time
	ObservedTimestamp    time.Time
	TraceId              *trace.TraceID
	SpanId               *trace.SpanID
	TraceFlags           *trace.TraceFlags
	SeverityText         *string
	SeverityNumber       *logs.SeverityNumber
	Body                 *string
	Resource             *resource.Resource
	InstrumentationScope *instrumentation.Scope
	Attributes           *[]attribute.KeyValue
}

// LogRecordStubFromReadableLogRecord returns a LogRecordStub populated from rl.
func LogRecordStubFromReadableLogRecord(rl logssdk.ReadableLogRecord) LogRecordStub {
	if rl == nil {
		return LogRecordStub{}
	}
	return LogRecordStub{
		Timestamp:            rl.Timestamp(),
		ObservedTimestamp:    rl.ObservedTimestamp(),
		TraceId:              rl.TraceId(),
		SpanId:               rl.SpanId(),
		TraceFlags:           rl.TraceFlags(),
		SeverityText:         rl.SeverityText(),
		SeverityNumber:       rl.SeverityNumber(),
		Body:                 rl.Body(),
		Resource:             rl.Resource(),
		InstrumentationScope: rl.InstrumentationScope(),
		Attributes:           rl.Attributes(),
	}
}

// Snapshot returns a read-only copy of the SpanStub.
func (s LogRecordStub) Snapshot() logssdk.ReadableLogRecord {
	return &logRecordSnapshot{
		timestamp:            s.Timestamp,
		observedTimestamp:    s.ObservedTimestamp,
		traceId:              s.TraceId,
		spanId:               s.SpanId,
		traceFlags:           s.TraceFlags,
		severityText:         s.SeverityText,
		severityNumber:       s.SeverityNumber,
		body:                 s.Body,
		resource:             s.Resource,
		instrumentationScope: s.InstrumentationScope,
		attributes:           s.Attributes,
	}
}

type logRecordSnapshot struct {
	// Embed the interface to implement the private method.
	logssdk.ReadableLogRecord
	timestamp            *time.Time
	observedTimestamp    time.Time
	traceId              *trace.TraceID
	spanId               *trace.SpanID
	traceFlags           *trace.TraceFlags
	severityText         *string
	severityNumber       *logs.SeverityNumber
	body                 *string
	resource             *resource.Resource
	instrumentationScope *instrumentation.Scope
	attributes           *[]attribute.KeyValue
}

func (r *logRecordSnapshot) Timestamp() *time.Time         { return r.timestamp }
func (r *logRecordSnapshot) ObservedTimestamp() time.Time  { return r.observedTimestamp }
func (r *logRecordSnapshot) TraceId() *trace.TraceID       { return r.traceId }
func (r *logRecordSnapshot) SpanId() *trace.SpanID         { return r.spanId }
func (r *logRecordSnapshot) TraceFlags() *trace.TraceFlags { return r.traceFlags }
func (r *logRecordSnapshot) InstrumentationScope() *instrumentation.Scope {
	return r.instrumentationScope
}
func (r *logRecordSnapshot) SeverityText() *string                { return r.severityText }
func (r *logRecordSnapshot) SeverityNumber() *logs.SeverityNumber { return r.severityNumber }
func (r *logRecordSnapshot) Body() *string                        { return r.body }
func (r *logRecordSnapshot) Resource() *resource.Resource         { return r.resource }
func (r *logRecordSnapshot) Attributes() *[]attribute.KeyValue    { return r.attributes }
