# OpenTelemetry Exporters

Once the OpenTelemetry SDK has created and processed telemetry, it needs to be exported. This package contains exporters
for this purpose.

## Exporter Packages

The following exporter packages are provided with the following OpenTelemetry signal support.

| Exporter Package                                                   | Logs |
|--------------------------------------------------------------------|------|
| github.com/kudarap/opentelemetry-logs-go/exporters/otlp/otlplogs | ✓    |
| github.com/kudarap/opentelemetry-logs-go/exporters/stdout        | ✓    |

## OTLP Logs exporter

Regarding [specification](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md#specify-protocol)
otlp exporter implemented `otlplogsgrpc` and `otlplogshttp` clients with next protocols
supported: `gprc`, `http/protobuf` and `http/json`.

If client is not configured explicitly ,The `OTEL_EXPORTER_OTLP_PROTOCOL`, `OTEL_EXPORTER_OTLP_LOGS_PROTOCOL`
environment variables specify the OTLP transport protocol. Supported values:

- `grpc` for protobuf-encoded data using gRPC wire format over HTTP/2 connection
- `http/protobuf` for protobuf-encoded data over HTTP connection
- `http/json` for JSON-encoded data over HTTP connection

```go
package main

func main() {
	ctx := context.Background()
	exporter, _ := otlplogs.NewExporter(ctx) // will create exporter with http client `http/protobuf` protocol by default

}

```

### OTEL client with grpc

Use `OTEL_EXPORTER_OTLP_PROTOCOL=grpc` env configuration or specify http client explicitly:

```go
exporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogsgrpc.NewClient()))
```

### OTEL client with http/protobuf

Use `OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf` env configuration or specify http client explicitly:

```go
exporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient(otlplogshttp.WithProtobufProtocol())))
```

### OTEL client with http/json

Use `OTEL_EXPORTER_OTLP_PROTOCOL=http/json` env configuration or specify http client explicitly:

```go
exporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient(otlplogshttp.WithJsonProtocol())))
```

## StdOut Logs exporter

The logging exporter prints the name of the log along with its attributes to stdout. It's mainly used for testing and
debugging.

```go
exporter, _ := stdoutlogs.NewExporter(ctx)
```

