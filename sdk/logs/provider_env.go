package logs

const (
	// see https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#exporter-selection
	logsExporterKey = "OTEL_LOGS_EXPORTER"

	logsExporterOTLP    = "otlp"
	logsExporterLogging = "logging"
)
