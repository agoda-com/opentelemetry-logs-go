# OpenTelemetry-Logs-Go

OpenTelemetry-Logs-Go is the [Go](https://golang.org) implementation of [OpenTelemetry](https://opentelemetry.io/) Logs.
It provides API to directly send logging data to observability platforms. It is an extension of official
[open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) to support Logs.

## Project Life Cycle

This project was created due log module freeze in
official [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) repository:

```
The Logs signal development is halted for this project while we stablize the Metrics SDK. 
No Logs Pull Requests are currently being accepted.
```

This project will be deprecated once official [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go)
repository Logs module will have status "Stable".

## Quick start

This is an implementation of [Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/) and not
intended to use by developers directly. It is provided for logging library authors to build log appenders, which use
this API to bridge between existing logging libraries and the OpenTelemetry log data model.

Example bellow will show how logging library could be instrumented with current API:

```go

const (
  instrumentationName = "otel/zap"
  instrumentationVersion = "0.0.1"
)

var (
  logger = otel.GetLoggerProvider().Logger(
      instrumentationName,
      logs.WithInstrumentationVersion(instrumentationVersion),
      logs.WithSchemaURL(semconv.SchemaURL),
  )
)

func (c otlpCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	
  lrc := logs.LogRecordConfig{
    Body: &ent.Message,
	...
  }
  logRecord := logs.NewLogRecord(cfg)
  logger.Emit(logRecord)
}
```

and application initialization code:

```go
package main

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go/otel"
	"github.com/agoda-com/opentelemetry-logs-go/otel/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/otel/exporters/otlp/otlplogs/otlplogshttp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	sdk "github.com/agoda-com/opentelemetry-logs-go/otel/sdk/logs"
)

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otlplogs-example"),
		semconv.ServiceVersion("0.0.1"),
	)
}

func main() {
	ctx := context.Background()

	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(otlplogs.New(ctx, otlplogshttp.NewClient())),
		sdk.WithResource(newResource()),
	)
	otel.SetLoggerProvider(loggerProvider)
	
	myInstrumentedLogger.Info("Hello OpenTelemetry")
}
```

