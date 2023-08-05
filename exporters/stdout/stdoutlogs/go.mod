module github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs

go 1.19

require (
	github.com/agoda-com/opentelemetry-logs-go v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/logs v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/sdk v0.1.1
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace (
	github.com/agoda-com/opentelemetry-logs-go => ../../..
	github.com/agoda-com/opentelemetry-logs-go/logs => ../../../logs
	github.com/agoda-com/opentelemetry-logs-go/sdk => ../../../sdk
)
