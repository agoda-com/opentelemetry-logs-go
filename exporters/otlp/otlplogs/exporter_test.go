package otlplogs_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/kudarap/opentelemetry-logs-go/sdk/logs/logstest"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
)

type client struct {
	uploadErr error
}

var _ otlplogs.Client = &client{}

func (c *client) Start(ctx context.Context) error {
	return nil
}

func (c *client) Stop(ctx context.Context) error {
	return nil
}

func (c *client) UploadLogs(ctx context.Context, protoLogs []*logspb.ResourceLogs) error {
	return c.uploadErr
}

func TestExporterClientError(t *testing.T) {
	ctx := context.Background()
	exp, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(&client{
		uploadErr: context.Canceled,
	}))
	assert.NoError(t, err)

	body := "Log record"
	logs := logstest.LogRecordStubs{{Body: &body}}.Snapshots()
	err = exp.Export(ctx, logs)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
	assert.True(t, strings.HasPrefix(err.Error(), "context canceled"), err)

	assert.NoError(t, exp.Shutdown(ctx))
}
