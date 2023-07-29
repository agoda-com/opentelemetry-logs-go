/*
Copyright Agoda Services Co.,Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package otlplogshttp

import (
	"context"
	"github.com/agoda-com/opentelemetry-logs-go"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"log"
	"time"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otlplogs-example"),
		semconv.ServiceVersion("0.0.1"),
	)
}

func doSomething() {

	logger := otel.GetLoggerProvider().Logger(
		instrumentationName,
		logs.WithInstrumentationVersion(instrumentationVersion),
		logs.WithSchemaURL(semconv.SchemaURL),
	)

	body := "Body"
	now := time.Now()
	cfg := logs.LogRecordConfig{
		Timestamp: &now,
		Body:      &body,
	}
	logRecord := logs.NewLogRecord(cfg)
	logger.Emit(logRecord)
}

func installExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	client := NewClient()
	exporter, _ := otlplogs.New(ctx, client)

	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}

func Example() {
	{
		ctx := context.Background()
		// Registers a tracer Provider globally.
		shutdown, err := installExportPipeline(ctx)
		if err != nil {
			log.Fatal(err)
		}
		doSomething()

		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
