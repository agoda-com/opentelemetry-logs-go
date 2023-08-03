module github.com/agoda-com/opentelemetry-logs-go/sdk

go 1.19

replace github.com/agoda-com/opentelemetry-logs-go => ../

require (
	github.com/agoda-com/opentelemetry-logs-go v0.0.1
	github.com/agoda-com/opentelemetry-logs-go/logs v0.0.1
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/agoda-com/opentelemetry-logs-go/logs => ../logs
