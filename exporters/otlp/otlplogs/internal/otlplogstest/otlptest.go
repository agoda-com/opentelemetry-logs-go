package otlplogstest

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	sdklogs "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

// RunEndToEndTest can be used by otlplogs.Client tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlplogs.Exporter, logsCollector LogsCollector) {
	pOpts := []sdklogs.LoggerProviderOption{
		sdklogs.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdklogs.WithBatchTimeout(5*time.Second),
			sdklogs.WithMaxExportBatchSize(10),
		),
	}
	tp1 := sdklogs.NewLoggerProvider(append(pOpts,
		sdklogs.WithResource(resource.NewSchemaless(
			attribute.String("rk1", "rv11)"),
			attribute.Int64("rk2", 5),
		)))...)

	tp2 := sdklogs.NewLoggerProvider(append(pOpts,
		sdklogs.WithResource(resource.NewSchemaless(
			attribute.String("rk1", "rv12)"),
			attribute.Float64("rk3", 6.5),
		)))...)

	tr1 := tp1.Logger("test-logger1")
	tr2 := tp2.Logger("test-logger2")
	// Now create few logs
	m := 4
	body := "TestLog"
	for i := 0; i < m; i++ {
		lr1 := logs.NewLogRecord(logs.LogRecordConfig{
			Body:       &body,
			Attributes: &[]attribute.KeyValue{attribute.Int64("i", int64(i))},
		})
		tr1.Emit(lr1)

		lr2 := logs.NewLogRecord(logs.LogRecordConfig{
			Body:       &body,
			Attributes: &[]attribute.KeyValue{attribute.Int64("i", int64(i))},
		})

		tr2.Emit(lr2)
	}

	func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := tp1.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a logger provider 1: %v", err)
		}
		if err := tp2.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a logger provider 2: %v", err)
		}
	}()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to stop the exporter: %v", err)
	}

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	if err := logsCollector.Stop(); err != nil {
		t.Fatalf("failed to stop the mock collector: %v", err)
	}

	// Now verify that we only got two resources
	rss := logsCollector.GetResourceLogs()
	if got, want := len(rss), 2; got != want {
		t.Fatalf("resource log count: got %d, want %d\n", got, want)
	}

	// Now verify logs and attributes for each resource log.
	for _, rs := range rss {
		if len(rs.ScopeLogs) == 0 {
			t.Fatalf("zero ScopeLogs")
		}
		if got, want := len(rs.ScopeLogs[0].LogRecords), m; got != want {
			t.Fatalf("log counts: got %d, want %d", got, want)
		}
		attrMap := map[int64]bool{}
		for _, s := range rs.ScopeLogs[0].LogRecords {
			if gotName, want := s.Body.GetStringValue(), "TestLog"; gotName != want {
				t.Fatalf("log name: got %s, want %s", gotName, want)
			}
			attrMap[s.Attributes[0].Value.Value.(*commonpb.AnyValue_IntValue).IntValue] = true
		}
		if got, want := len(attrMap), m; got != want {
			t.Fatalf("log attribute unique values: got %d  want %d", got, want)
		}
		for i := 0; i < m; i++ {
			_, ok := attrMap[int64(i)]
			if !ok {
				t.Fatalf("log with attribute %d missing", i)
			}
		}
	}
}
