package otlplogs

import (
	"context"
	"errors"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/logstransform"
	logssdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"sync"
)

var (
	errAlreadyStarted = errors.New("already started")
)

type Exporter struct {
	client Client

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once
}

// Start establishes a connection to the receiving endpoint.
func (e *Exporter) Start(ctx context.Context) error {
	var err = errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()
		err = e.client.Start(ctx)
	})

	return err
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
func New(ctx context.Context, client Client) (*Exporter, error) {
	exp := &Exporter{
		client: client,
	}
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}
