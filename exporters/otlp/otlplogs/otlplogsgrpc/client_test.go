package otlplogsgrpc_test

import (
	"context"
	"fmt"
	"github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/otlplogstest"
	"github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	"github.com/kudarap/opentelemetry-logs-go/logs"
	sdklogs "github.com/kudarap/opentelemetry-logs-go/sdk/logs"
	"github.com/kudarap/opentelemetry-logs-go/sdk/logs/logstest"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

var body = "Log Record 0"
var roLogRecords = logstest.LogRecordStubs{{Body: &body}}.Snapshots()

func contextWithTimeout(parent context.Context, t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	d, ok := t.Deadline()
	if !ok {
		d = time.Now().Add(timeout)
	} else {
		d = d.Add(-1 * time.Millisecond)
		now := time.Now()
		if d.Sub(now) > timeout {
			d = now.Add(timeout)
		}
	}
	return context.WithDeadline(parent, d)
}

func TestNewEndToEnd(t *testing.T) {
	tests := []struct {
		name           string
		additionalOpts []otlplogsgrpc.Option
	}{
		{
			name: "StandardExporter",
		},
		{
			name: "WithCompressor",
			additionalOpts: []otlplogsgrpc.Option{
				otlplogsgrpc.WithCompressor(gzip.Name),
			},
		},
		{
			name: "WithServiceConfig",
			additionalOpts: []otlplogsgrpc.Option{
				otlplogsgrpc.WithServiceConfig("{}"),
			},
		},
		{
			name: "WithDialOptions",
			additionalOpts: []otlplogsgrpc.Option{
				otlplogsgrpc.WithDialOption(grpc.WithBlock()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newExporterEndToEndTest(t, test.additionalOpts)
		})
	}
}

func newGRPCExporter(t *testing.T, ctx context.Context, endpoint string, additionalOpts ...otlplogsgrpc.Option) *otlplogs.Exporter {
	opts := []otlplogsgrpc.Option{
		otlplogsgrpc.WithInsecure(),
		otlplogsgrpc.WithEndpoint(endpoint),
		otlplogsgrpc.WithReconnectionPeriod(50 * time.Millisecond),
	}

	opts = append(opts, additionalOpts...)
	client := otlplogsgrpc.NewClient(opts...)
	exp, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(client))
	if err != nil {
		t.Fatalf("failed to create a new collector exporter: %v", err)
	}
	return exp
}

func newExporterEndToEndTest(t *testing.T, additionalOpts []otlplogsgrpc.Option) {
	mc := runMockCollector(t)

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint, additionalOpts...)
	t.Cleanup(func() {
		ctx, cancel := contextWithTimeout(ctx, t, 10*time.Second)
		defer cancel()

		require.NoError(t, exp.Shutdown(ctx))
	})

	// RunEndToEndTest closes mc.
	otlplogstest.RunEndToEndTest(ctx, t, exp, mc)
}

func TestExporterShutdown(t *testing.T) {
	mc := runMockCollectorAtEndpoint(t, "localhost:0")
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	factory := func() otlplogs.Client {
		return otlplogsgrpc.NewClient(
			otlplogsgrpc.WithEndpoint(mc.endpoint),
			otlplogsgrpc.WithInsecure(),
			otlplogsgrpc.WithDialOption(grpc.WithBlock()),
		)
	}
	otlplogstest.RunExporterShutdownTest(t, factory)
}

func TestNewInvokeStartThenStopManyTimes(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	// Invoke Start numerous times, should return errAlreadyStarted
	for i := 0; i < 10; i++ {
		if err := exp.Start(ctx); err == nil || !strings.Contains(err.Error(), "already started") {
			t.Fatalf("#%d unexpected Start error: %v", i, err)
		}
	}

	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to Shutdown the exporter: %v", err)
	}
	// Invoke Shutdown numerous times
	for i := 0; i < 10; i++ {
		if err := exp.Shutdown(ctx); err != nil {
			t.Fatalf(`#%d got error (%v) expected none`, i, err)
		}
	}
}

// This test takes a long time to run: to skip it, run tests using: -short.
func TestNewCollectorOnBadConnection(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping this long running test")
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to grab an available port: %v", err)
	}
	// Firstly close the "collector's" channel: optimistically this endpoint won't get reused ASAP
	// However, our goal of closing it is to simulate an unavailable connection
	_ = ln.Close()

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	endpoint := fmt.Sprintf("localhost:%s", collectorPortStr)
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNewWithEndpoint(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	_ = exp.Shutdown(ctx)
}

func TestNewWithHeaders(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlplogsgrpc.WithHeaders(map[string]string{"header1": "value1"}))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.Export(ctx, roLogRecords))

	headers := mc.getHeaders()
	require.Regexp(t, "OTel OTLP Exporter Go/1\\..*", headers.Get("user-agent"))
	require.Len(t, headers.Get("header1"), 1)
	assert.Equal(t, "value1", headers.Get("header1")[0])
}

//func TestExportLogsTimeoutHonored(t *testing.T) {
//	ctx, cancel := contextWithTimeout(context.Background(), t, 1*time.Minute)
//	t.Cleanup(cancel)
//
//	mc := runMockCollector(t)
//	exportBlock := make(chan struct{})
//	mc.logsSvc.exportBlock = exportBlock
//	t.Cleanup(func() { require.NoError(t, mc.stop()) })
//
//	exp := newGRPCExporter(
//		t,
//		ctx,
//		mc.endpoint,
//		otlplogsgrpc.WithTimeout(1*time.Nanosecond),
//		otlplogsgrpc.WithRetry(otlplogsgrpc.RetryConfig{Enabled: false}),
//	)
//	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
//
//	err := exp.Export(ctx, roLogRecords)
//	// Release the export so everything is cleaned up on shutdown.
//	close(exportBlock)
//
//	unwrapped := errors.Unwrap(err)
//	require.Equal(t, codes.DeadlineExceeded, status.Convert(unwrapped).Code())
//	require.True(t, strings.HasPrefix(err.Error(), "logRecords export: "), err)
//}

func TestNewWithMultipleAttributeTypes(t *testing.T) {
	mc := runMockCollector(t)

	ctx, cancel := contextWithTimeout(context.Background(), t, 10*time.Second)
	t.Cleanup(cancel)

	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	tp := sdklogs.NewLoggerProvider(
		sdklogs.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdklogs.WithBatchTimeout(5*time.Second),
			sdklogs.WithMaxExportBatchSize(10),
		),
	)
	t.Cleanup(func() { require.NoError(t, tp.Shutdown(ctx)) })

	tr := tp.Logger("test-logger")
	testKvs := []attribute.KeyValue{
		attribute.Int("Int", 1),
		attribute.Int64("Int64", int64(3)),
		attribute.Float64("Float64", 2.22),
		attribute.Bool("Bool", true),
		attribute.String("String", "test"),
	}
	lrc := logs.LogRecordConfig{
		Attributes: &testKvs,
	}

	logRecord := logs.NewLogRecord(lrc)

	tr.Emit(logRecord)

	// Flush and close.
	func() {
		ctx, cancel := contextWithTimeout(ctx, t, 10*time.Second)
		defer cancel()
		require.NoError(t, tp.Shutdown(ctx))
	}()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	require.NoError(t, exp.Shutdown(ctx))

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	require.NoError(t, mc.stop())

	// Now verify that we only got one logs
	rss := mc.getLogRecords()
	if got, want := len(rss), 1; got != want {
		t.Fatalf("resource logs count: got %d, want %d\n", got, want)
	}

	expected := []*commonpb.KeyValue{
		{
			Key: "Int",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 1,
				},
			},
		},
		{
			Key: "Int64",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 3,
				},
			},
		},
		{
			Key: "Float64",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_DoubleValue{
					DoubleValue: 2.22,
				},
			},
		},
		{
			Key: "Bool",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_BoolValue{
					BoolValue: true,
				},
			},
		},
		{
			Key: "String",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: "test",
				},
			},
		},
	}

	// Verify attributes
	if !assert.Len(t, rss[0].Attributes, len(expected)) {
		t.Fatalf("attributes count: got %d, want %d\n", len(rss[0].Attributes), len(expected))
	}
	for i, actual := range rss[0].Attributes {
		if a, ok := actual.Value.Value.(*commonpb.AnyValue_DoubleValue); ok {
			e, ok := expected[i].Value.Value.(*commonpb.AnyValue_DoubleValue)
			if !ok {
				t.Errorf("expected AnyValue_DoubleValue, got %T", expected[i].Value.Value)
				continue
			}
			if !assert.InDelta(t, e.DoubleValue, a.DoubleValue, 0.01) {
				continue
			}
			e.DoubleValue = a.DoubleValue
		}
		assert.Equal(t, expected[i], actual)
	}
}

func TestStartErrorInvalidAddress(t *testing.T) {
	client := otlplogsgrpc.NewClient(
		otlplogsgrpc.WithInsecure(),
		// Validate the connection in Start (which should return the error).
		otlplogsgrpc.WithDialOption(
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		),
		otlplogsgrpc.WithEndpoint("invalid"),
		otlplogsgrpc.WithReconnectionPeriod(time.Hour),
	)
	err := client.Start(context.Background())
	assert.EqualError(t, err, `connection error: desc = "transport: error while dialing: dial tcp: address invalid: missing port in address"`)
}

func TestEmptyData(t *testing.T) {
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })

	assert.NoError(t, exp.Export(ctx, nil))
}

func TestPartialSuccess(t *testing.T) {
	mc := runMockCollectorWithConfig(t, &mockConfig{
		partial: &collogspb.ExportLogsPartialSuccess{
			RejectedLogRecords: 2,
			ErrorMessage:       "partially successful",
		},
	})
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	errs := []error{}
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errs = append(errs, err)
	}))
	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint)
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.Export(ctx, roLogRecords))

	require.Equal(t, 1, len(errs))
	require.Contains(t, errs[0].Error(), "partially successful")
	require.Contains(t, errs[0].Error(), "2 logs rejected")
}

func TestCustomUserAgent(t *testing.T) {
	customUserAgent := "custom-user-agent"
	mc := runMockCollector(t)
	t.Cleanup(func() { require.NoError(t, mc.stop()) })

	ctx := context.Background()
	exp := newGRPCExporter(t, ctx, mc.endpoint,
		otlplogsgrpc.WithDialOption(grpc.WithUserAgent(customUserAgent)))
	t.Cleanup(func() { require.NoError(t, exp.Shutdown(ctx)) })
	require.NoError(t, exp.Export(ctx, roLogRecords))

	headers := mc.getHeaders()
	require.Contains(t, headers.Get("user-agent")[0], customUserAgent)
}
