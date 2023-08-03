package otlplogs

import (
	//	"context"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/internal/otlpconfig"
	//	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttpjson"
)

var (
	// Clients TODO: make private
	Clients = map[otlpconfig.Protocol]Client{}
)

func init() {
	//	c, _ := otlplogshttpjson.New(context.Background())
	//	println(c)
}
