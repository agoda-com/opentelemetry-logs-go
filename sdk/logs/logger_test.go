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
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"testing"
	"time"
)

func TestLogsReadWriteAPIFormat(t *testing.T) {

	traceID, _ := trace.TraceIDFromHex("80f198ee56343ba864fe8b2a57d3eff7")
	spanID, _ := trace.SpanIDFromHex("2a00000000000000")

	spanCtxCfg := trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
		Remote:     false,
	}

	spanCtx := trace.NewSpanContext(spanCtxCfg)

	body := "My Log Message"
	severityText := "INFO"
	severityNumber := logs.INFO
	resource := sdkresource.NewWithAttributes("http", semconv.HTTPURL("testurl"))
	instrumentationScope := instrumentation.Scope{
		Name:      "test",
		Version:   "testVersion",
		SchemaURL: "http",
	}
	attributes := []attribute.KeyValue{semconv.HTTPURL("testurl")}
	timestamp := time.Now()

	record := newReadWriteLogRecord(
		&spanCtx,
		&body,
		&severityText,
		&severityNumber,
		resource,
		&instrumentationScope,
		&attributes,
		&timestamp,
	)

	assert.Equal(t, "My Log Message", *record.Body())

}
