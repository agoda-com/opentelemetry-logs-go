module github.com/agoda-com/opentelemetry-logs-go

go 1.20

require (
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/golang/protobuf v1.5.3
	github.com/stretchr/testify v1.8.3
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	go.opentelemetry.io/proto/otlp v0.19.0
	go.uber.org/goleak v1.2.1
	google.golang.org/grpc v1.55.0
)

require google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
