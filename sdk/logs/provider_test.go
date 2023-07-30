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

package logs

import (
	//	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"testing"
)

const (
	instrumentationName    = "otel/zap"
	instrumentationVersion = "0.0.1"
)

func TestLogsProvider(t *testing.T) {

	noopProvider := logs.NewNoopLoggerProvider()

	noopLogger := noopProvider.Logger(instrumentationName,
		logs.WithInstrumentationVersion(instrumentationVersion),
		logs.WithSchemaURL(semconv.SchemaURL),
	)

	logRecordExporter := NewTestExporter()

	batchOtlpLogger := NewLoggerProvider(
		WithLogRecordProcessor(NewBatchLogRecordProcessor(logRecordExporter)),
		WithResource(
			resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("unit_test"))),
	).Logger("otlp",
		logs.WithInstrumentationVersion(instrumentationName),
		logs.WithSchemaURL(instrumentationVersion),
		logs.WithInstrumentationAttributes(semconv.HostName("some.host")),
	)

	body := "body"

	logRecord := logs.NewLogRecord(logs.LogRecordConfig{
		Body: &body,
	})

	noopLogger.Emit(logRecord)
	batchOtlpLogger.Emit(logRecord)

}
