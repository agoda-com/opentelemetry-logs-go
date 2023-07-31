package logs

const (
	// see https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#exporter-selection
	logsExporterKey = "OTEL_LOGS_EXPORTER"

	logsExporterOTLP    = "otlp"
	logsExporterLogging = "logging"

	// TODO: move to exporters/otlp package
	// see https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_protocol
	exporterOtlpProtocolKey     = "OTEL_EXPORTER_OTLP_PROTOCOL"
	logsExporterOtlpProtocolKey = "OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"

	logsExporterOtlpProtocolGrpc         = "grpc"
	logsExporterOtlpProtocolHttpProtobuf = "http/protobuf"
	logsExporterOtlpProtocolHttpJson     = "http/json"
)
