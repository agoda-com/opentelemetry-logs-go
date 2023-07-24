package otlplogshttp

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
)

// New constructs a new Exporter and starts it.
func New(ctx context.Context) (*otlplogs.Exporter, error) {
	return otlplogs.New(ctx, NewClient()), nil
}
