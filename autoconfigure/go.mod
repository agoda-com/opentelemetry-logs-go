module github.com/agoda-com/opentelemetry-logs-go/autoconfigure

go 1.19

require (
	github.com/agoda-com/opentelemetry-logs-go v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/logs v0.1.1
	github.com/agoda-com/opentelemetry-logs-go/sdk v0.1.1
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
)

require (
	github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/internal/retry v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/agoda-com/opentelemetry-logs-go => ../
	github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/internal/retry => ../exporters/otlp/internal/retry
	github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs => ../exporters/otlp/otlplogs
	github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs => ../exporters/stdout/stdoutlogs
	github.com/agoda-com/opentelemetry-logs-go/logs => ../logs
	github.com/agoda-com/opentelemetry-logs-go/sdk => ../sdk
)
