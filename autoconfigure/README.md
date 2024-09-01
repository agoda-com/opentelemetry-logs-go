# OpenTelemetry SDK Autoconfigure

This package is `Experimental` and subject to change/remove in the future

This artifact implements environment-based autoconfiguration of the OpenTelemetry SDK. This can be
an alternative to programmatic configuration using the normal SDK.

All options support being passed as environment variables, e.g., `OTEL_LOGS_EXPORTER=otlp`.

```golang
// export OTEL_SERVICE_NAME=otlplogs-example
// export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
// export OTEL_LOGS_EXPORTER=otlp,logging
// export OTEL_RESOURCE_ATTRIBUTES="key1=value1,key2=value2"

package main

import (
	"os"
	"context"
	"go.opentelemetry.io/otel/sdk/resource"
  "github.com/kudarap/opentelemetry-logs-go/autoconfigure/sdk/logs"
)

func main() {

	// It is possible to mix with explicit configuration
	host, _ := os.Hostname()
	attrs := resource.NewWithAttributes(
		semconv.HostName(host),
	)

	ctx := context.Background()
	loggerProvider := autosdk.NewLoggerProvider(ctx, autosdk.WithResource(attrs))
	defer func() {
		loggerProvider.Shutdown(ctx)
	}()
}

```

## Contents

- [General Configuration](#general-configuration)
    * [Exporters](#exporters)
        + [OTLP exporter (log exporters)](#otlp-exporter-log-exporters)
            - [OTLP exporter retry](#otlp-exporter-retry)
        + [Logging exporter](#logging-exporter)
    * [OpenTelemetry Resource](#opentelemetry-resource)
    * [Attribute limits](#attribute-limits)
- [Logger provider](#logger-provider)
- [Batch log record processor](#batch-log-record-processor)
- [Customizing the OpenTelemetry SDK](#customizing-the-opentelemetry-sdk)

## General configuration

### Exporters

The following configuration properties are common to all exporters:

| Environment variable | Purpose                                                                                                                    |
|----------------------|----------------------------------------------------------------------------------------------------------------------------|
| OTEL_LOGS_EXPORTER   | List of exporters to be used for logging, separated by commas. Default is `otlp`. `none` means no autoconfigured exporter. |

#### OTLP exporter (log exporters)

The [OpenTelemetry Protocol Exporter](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md)
span, metric, and log exporters

| Environment variable                       | Description                                                                                                                                                                                                                                                                                                                                                                            |
|--------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| OTEL_LOGS_EXPORTER                         | Select the OpenTelemetry exporter for logs (default `otlp`)                                                                                                                                                                                                                                                                                                                            |
| OTEL_EXPORTER_OTLP_ENDPOINT                | The OTLP traces, metrics, and logs endpoint to connect to. Must be a URL with a scheme of either `http` or `https` based on the use of TLS. If protocol is `http/protobuf` the version and signal will be appended to the path (e.g. `v1/logs`). Default is `http://localhost:4317` when protocol is `grpc`, and `http://localhost:4318/v1/{signal}` when protocol is `http/protobuf`. |
| OTEL_EXPORTER_OTLP_LOGS_ENDPOINT           | The OTLP logs endpoint to connect to. Must be a URL with a scheme of either `http` or `https` based on the use of TLS. Default is `http://localhost:4317` when protocol is `grpc`, and `http://localhost:4318/v1/logs` when protocol is `http/protobuf`.                                                                                                                               |
| OTEL_EXPORTER_OTLP_PROTOCOL                | The transport protocol to use on OTLP trace, metric, and log requests. Options include `grpc` and `http/protobuf`. Default is `grpc`.                                                                                                                                                                                                                                                  |
| OTEL_EXPORTER_OTLP_LOGS_PROTOCOL           | The transport protocol to use on OTLP log requests. Options include `grpc` and `http/protobuf`. Default is `grpc`.                                                                                                                                                                                                                                                                     |
| OTEL_EXPORTER_OTLP_CERTIFICATE             | The path to the file containing trusted certificates to use when verifying an OTLP trace, metric, or log server's TLS credentials. The file should contain one or more X.509 certificates in PEM format. By default the host platform's trusted root certificates are used.                                                                                                            |
| OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE        | The path to the file containing trusted certificates to use when verifying an OTLP log server's TLS credentials. The file should contain one or more X.509 certificates in PEM format. By default the host platform's trusted root certificates are used.                                                                                                                              |
| OTEL_EXPORTER_OTLP_CLIENT_KEY              | The path to the file containing private client key to use when verifying an OTLP trace, metric, or log client's TLS credentials. The file should contain one private key PKCS8 PEM format. By default no client key is used.                                                                                                                                                           |
| OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY         | The path to the file containing private client key to use when verifying an OTLP log client's TLS credentials. The file should contain one private key PKCS8 PEM format. By default no client key file is used.                                                                                                                                                                        |
| OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE      | The path to the file containing trusted certificates to use when verifying an OTLP trace, metric, or log client's TLS credentials. The file should contain one or more X.509 certificates in PEM format. By default no chain file is used.                                                                                                                                             |
| OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE | The path to the file containing trusted certificates to use when verifying an OTLP log server's TLS credentials. The file should contain one or more X.509 certificates in PEM format. By default no chain file is used.                                                                                                                                                               |
| OTEL_EXPORTER_OTLP_INSECURE                | Whether to enable trace, metric or log client's transport security for the exporter's gRPC connection. This option only applies to OTLP/gRPC when an endpoint is provided without the http or https scheme. Default `false`                                                                                                                                                            |
| OTEL_EXPORTER_OTLP_LOGS_INSECURE           | Whether to enable log client's transport security for the exporter's gRPC connection. This option only applies to OTLP/gRPC when an endpoint is provided without the http or https scheme. Default `false`                                                                                                                                                                             |
| OTEL_EXPORTER_OTLP_HEADERS                 | Key-value pairs separated by commas to pass as request headers on OTLP trace, metric, and log requests.                                                                                                                                                                                                                                                                                |
| OTEL_EXPORTER_OTLP_LOGS_HEADERS            | Key-value pairs separated by commas to pass as request headers on OTLP logs requests.                                                                                                                                                                                                                                                                                                  |
| OTEL_EXPORTER_OTLP_COMPRESSION             | The compression type to use on OTLP trace, metric, and log requests. Options include `gzip`. By default no compression will be used.                                                                                                                                                                                                                                                   |
| OTEL_EXPORTER_OTLP_LOGS_COMPRESSION        | The compression type to use on OTLP log requests. Options include `gzip`. By default no compression will be used.                                                                                                                                                                                                                                                                      |
| OTEL_EXPORTER_OTLP_TIMEOUT                 | The maximum waiting time, in milliseconds, allowed to send each OTLP trace, metric, and log batch. Default is `10000`.                                                                                                                                                                                                                                                                 |
| OTEL_EXPORTER_OTLP_LOGS_TIMEOUT            | The maximum waiting time, in milliseconds, allowed to send each OTLP log batch. Default is `10000`.                                                                                                                                                                                                                                                                                    |

To configure the service name for the OTLP exporter, add the `service.name` key
to the OpenTelemetry Resource ([see below](#opentelemetry-resource)),
e.g. `OTEL_RESOURCE_ATTRIBUTES=service.name=myservice`.

##### OTLP exporter retry

[OTLP](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md#otlpgrpc-response)
requires
that [transient](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md#retry)
errors be handled with a retry strategy. When retry is enabled, retryable gRPC status codes will be retried using an
exponential backoff with jitter algorithm as described in
the [gRPC Retry Design](https://github.com/grpc/proposal/blob/master/A6-client-retries.md#exponential-backoff).

The policy has the following configuration, which there is currently no way to customize.

- `maxAttempts`: The maximum number of attempts, including the original request. Defaults to `5`.
- `initialBackoff`: The initial backoff duration. Defaults to `1s`
- `maxBackoff`: The maximum backoff duration. Defaults to `5s`.
- `backoffMultiplier` THe backoff multiplier. Defaults to `1.5`.

#### Logging exporter

The logging exporter prints the name of the span along with its attributes to stdout. It's mainly used for testing and
debugging.

| Environment variable       | Description                          |
|----------------------------|--------------------------------------|
| OTEL_LOGS_EXPORTER=logging | Select the logging exporter for logs |

### OpenTelemetry Resource

The [OpenTelemetry Resource](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/resource/sdk.md)
is a representation of the entity producing telemetry.

| Environment variable     | Description                                                                                                |
|--------------------------|------------------------------------------------------------------------------------------------------------|
| OTEL_RESOURCE_ATTRIBUTES | Specify resource attributes in the following format: key1=val1,key2=val2,key3=val3                         |
| OTEL_SERVICE_NAME        | Specify logical service name. Takes precedence over `service.name` defined with `OTEL_RESOURCE_ATTRIBUTES` |

You almost always want to specify
the [`service.name`](https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/resource/semantic_conventions#service)
for your application.
It corresponds to how you describe the application, for example `authservice` could be an application that authenticates
requests, and `cats` could be an application that returns information about [cats](https://en.wikipedia.org/wiki/Cat).
You would specify that by setting service name property in one of the following ways:

* directly via `OTEL_SERVICE_NAME=authservice`
* by `service.name` resource attribute like `OTEL_RESOURCE_ATTRIBUTES=service.name=authservice`,

If not specified, SDK defaults the service name to `unknown_service:path_to_class`.

### Attribute limits

These properties can be used to control the maximum number and length of attributes.

| Environment variable              | Description                                                           |
|-----------------------------------|-----------------------------------------------------------------------|
| OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT | The maximum length of attribute values. By default there is no limit. |
| OTEL_ATTRIBUTE_COUNT_LIMIT        | The maximum number of attributes. Default is `128`.                   |

## Logger provider

The following configuration options are specific to SDK `LoggerProvider`.
See [general configuration](#general-configuration) for general configuration.

## Batch log record processor

| Environment variable            | Description                                                                        |
|---------------------------------|------------------------------------------------------------------------------------|
| OTEL_BLRP_SCHEDULE_DELAY        | The interval, in milliseconds, between two consecutive exports. Default is `1000`. |
| OTEL_BLRP_MAX_QUEUE_SIZE        | The maximum queue size. Default is `2048`.                                         |
| OTEL_BLRP_MAX_EXPORT_BATCH_SIZE | The maximum batch size. Default is `512`.                                          |
| OTEL_BLRP_EXPORT_TIMEOUT        | The maximum allowed time, in milliseconds, to export data. Default is `30000`.     |

## Customizing the OpenTelemetry SDK

Autoconfiguration currently doesn't expose any methods for customization. Only default `otlp` and `logging` providers
can be configured. This is subject to improve in the future
