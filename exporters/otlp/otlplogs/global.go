package otlplogs

import (
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/otlpconfig"
)

var (
	// Clients TODO: make private
	Clients = map[otlpconfig.Protocol]Client{}
)
