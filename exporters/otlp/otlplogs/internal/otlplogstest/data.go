package otlplogstest

import (
	"github.com/kudarap/opentelemetry-logs-go/logs"
	logssdk "github.com/kudarap/opentelemetry-logs-go/sdk/logs"
	"github.com/kudarap/opentelemetry-logs-go/sdk/logs/logstest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func SingleReadableLogRecord() []logssdk.ReadableLogRecord {

	time := time.Now()
	tid := trace.TraceID{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5, 6, 7, 8, 9}
	sid := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	tf := trace.FlagsSampled
	body := "TestMessage"
	is := instrumentation.Scope{
		Name:    "bar",
		Version: "0.0.0",
	}
	st := "INFO"
	sn := logs.INFO

	return logstest.LogRecordStubs{
		logstest.LogRecordStub{
			Timestamp:            &time,
			ObservedTimestamp:    time,
			TraceId:              &tid,
			SpanId:               &sid,
			TraceFlags:           &tf,
			SeverityText:         &st,
			SeverityNumber:       &sn,
			Body:                 &body,
			Resource:             resource.NewSchemaless(attribute.String("a", "b")),
			InstrumentationScope: &is,
			Attributes:           &[]attribute.KeyValue{},
		},
	}.Snapshots()
}
