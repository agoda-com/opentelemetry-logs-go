package logs

import (
	"os"
	"strings"
)

const (
	// see https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#exporter-selection
	logsExporterKey = "OTEL_LOGS_EXPORTER"

	logsExporterNone    = "none"
	logsExporterOTLP    = "otlp"
	logsExporterLogging = "logging"
)

func exportersFromEnv() ([]string, bool) {
	exportersEnv, defined := os.LookupEnv(logsExporterKey)
	exporters := strings.Split(exportersEnv, ",")
	return exporters, defined
}
