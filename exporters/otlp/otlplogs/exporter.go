package otlplogs

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/logstransform"
	logssdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"sync"
)

type Exporter struct {
	client Client

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once
}

func (e *Exporter) Shutdown(ctx context.Context) error {
	e.mu.RLock()
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	var err error

	e.stopOnce.Do(func() {
		err = e.client.Stop(ctx)
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
	})

	return err
}

// Export exports a batch of logs.
func (e *Exporter) Export(ctx context.Context, ll []logssdk.ReadableLogRecord) error {
	protoLogs := logstransform.Logs(ll)
	if len(protoLogs) == 0 {
		return nil
	}

	err := e.client.UploadLogs(ctx, protoLogs)
	if err != nil {
		return err
	}
	return nil
}

// New creates new Exporter with provided client
func New(_ context.Context, client Client) *Exporter {
	return &Exporter{
		client: client,
	}
}
