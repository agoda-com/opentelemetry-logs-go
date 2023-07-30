package otlplogsgrpc_test

import (
	"context"
	"fmt"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/otlplogstest"
	"github.com/stretchr/testify/require"
	collectorlogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net"
	"sync"
	"testing"
)

func makeMockCollector(t *testing.T, mockConfig *mockConfig) *mockCollector {
	return &mockCollector{
		t: t,
		logsSvc: &mockLogsService{
			storage: otlplogstest.NewLogsStorage(),
			errors:  mockConfig.errors,
			partial: mockConfig.partial,
		},
		stopped: make(chan struct{}),
	}
}

type mockLogsService struct {
	collectorlogspb.UnimplementedLogsServiceServer

	errors      []error
	partial     *collectorlogspb.ExportLogsPartialSuccess
	requests    int
	mu          sync.RWMutex
	storage     otlplogstest.LogsStorage
	headers     metadata.MD
	exportBlock chan struct{}
}

func (mts *mockLogsService) getHeaders() metadata.MD {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.headers
}

func (mts *mockLogsService) getLogs() []*logspb.LogRecord {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.storage.GetLogRecords()
}

func (mts *mockLogsService) getResourceLogs() []*logspb.ResourceLogs {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.storage.GetResourceLogs()
}

func (mts *mockLogsService) Export(ctx context.Context, exp *collectorlogspb.ExportLogsServiceRequest) (*collectorlogspb.ExportLogsServiceResponse, error) {
	mts.mu.Lock()
	defer func() {
		mts.requests++
		mts.mu.Unlock()
	}()

	if mts.exportBlock != nil {
		// Do this with the lock held so the mockCollector.Stop does not
		// abandon cleaning up resources.
		<-mts.exportBlock
	}

	reply := &collectorlogspb.ExportLogsServiceResponse{
		PartialSuccess: mts.partial,
	}
	if mts.requests < len(mts.errors) {
		idx := mts.requests
		return reply, mts.errors[idx]
	}

	mts.headers, _ = metadata.FromIncomingContext(ctx)
	mts.storage.AddLogs(exp)
	return reply, nil
}

type mockCollector struct {
	t *testing.T

	logsSvc *mockLogsService

	endpoint string
	stopFunc func()
	stopOnce sync.Once
	stopped  chan struct{}
}

type mockConfig struct {
	errors   []error
	endpoint string
	partial  *collectorlogspb.ExportLogsPartialSuccess
}

var _ collectorlogspb.LogsServiceServer = (*mockLogsService)(nil)

var errAlreadyStopped = fmt.Errorf("already stopped")

func (mc *mockCollector) stop() error {
	err := errAlreadyStopped
	mc.stopOnce.Do(func() {
		err = nil
		if mc.stopFunc != nil {
			mc.stopFunc()
		}
	})
	// Wait until gRPC server is down.
	<-mc.stopped

	// Getting the lock ensures the logsSvc is done flushing.
	mc.logsSvc.mu.Lock()
	defer mc.logsSvc.mu.Unlock()

	return err
}

func (mc *mockCollector) Stop() error {
	return mc.stop()
}

func (mc *mockCollector) getLogRecords() []*logspb.LogRecord {
	return mc.logsSvc.getLogs()
}

func (mc *mockCollector) getResourceLogs() []*logspb.ResourceLogs {
	return mc.logsSvc.getResourceLogs()
}

func (mc *mockCollector) GetResourceLogs() []*logspb.ResourceLogs {
	return mc.getResourceLogs()
}

func (mc *mockCollector) getHeaders() metadata.MD {
	return mc.logsSvc.getHeaders()
}

// runMockCollector is a helper function to create a mock Collector.
func runMockCollector(t *testing.T) *mockCollector {
	t.Helper()
	return runMockCollectorAtEndpoint(t, "localhost:0")
}

func runMockCollectorAtEndpoint(t *testing.T, endpoint string) *mockCollector {
	t.Helper()
	return runMockCollectorWithConfig(t, &mockConfig{endpoint: endpoint})
}

func runMockCollectorWithConfig(t *testing.T, mockConfig *mockConfig) *mockCollector {
	t.Helper()
	ln, err := net.Listen("tcp", mockConfig.endpoint)
	require.NoError(t, err, "net.Listen")

	srv := grpc.NewServer()
	mc := makeMockCollector(t, mockConfig)
	collectorlogspb.RegisterLogsServiceServer(srv, mc.logsSvc)
	go func() {
		_ = srv.Serve(ln)
		close(mc.stopped)
	}()

	mc.endpoint = ln.Addr().String()
	mc.stopFunc = srv.Stop

	// Wait until gRPC server is up.
	conn, err := grpc.Dial(mc.endpoint, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "grpc.Dial")
	require.NoError(t, conn.Close(), "conn.Close")

	return mc
}
