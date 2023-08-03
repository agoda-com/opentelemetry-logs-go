# OpenTelemetry Exporters

Once the OpenTelemetry SDK has created and processed telemetry, it needs to be exported. This package contains exporters
for this purpose.

## Exporter Packages

The following exporter packages are provided with the following OpenTelemetry signal support.

| Exporter Package                                                   | Logs |
|--------------------------------------------------------------------|------|
| github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs | ✓    |
| github.com/agoda-com/opentelemetry-logs-go/exporters/stdout        | ✓    |
