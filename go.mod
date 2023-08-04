module github.com/agoda-com/opentelemetry-logs-go

go 1.19

require (
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.16.0
	github.com/agoda-com/opentelemetry-logs-go/logs v0.0.1
)

require (
	github.com/kr/pretty v0.1.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/agoda-com/opentelemetry-logs-go/logs => ./logs
