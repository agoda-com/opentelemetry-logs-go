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

package logstransform

import (
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"time"
)

// Logs transforms OpenTelemetry LogRecord's into a OTLP ResourceLogs
func Logs(sdl []sdk.ReadableLogRecord) []*logspb.ResourceLogs {

	var resourceLogs []*logspb.ResourceLogs

	for _, sd := range sdl {

		var body *commonpb.AnyValue = nil
		if sd.Body() != nil {
			body = &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: *sd.Body(),
				},
			}
		}

		var traceIDBytes []byte
		if sd.TraceId() != nil {
			tid := *sd.TraceId()
			traceIDBytes = tid[:]
		}
		var spanIDBytes []byte
		if sd.SpanId() != nil {
			sid := *sd.SpanId()
			spanIDBytes = sid[:]
		}
		var traceFlags byte = 0
		if sd.TraceFlags() != nil {
			tf := *sd.TraceFlags()
			traceFlags = byte(tf)
		}
		var ts time.Time
		if sd.Timestamp() != nil {
			ts = *sd.Timestamp()
		} else {
			ts = sd.ObservedTimestamp()
		}

		var kv []*commonpb.KeyValue
		if sd.Attributes() != nil {
			kv = KeyValues(*sd.Attributes())
		}

		var st = ""
		if sd.SeverityText() != nil {
			st = *sd.SeverityText()
		}

		var sn = logspb.SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
		if sd.SeverityNumber() != nil {
			sn = logspb.SeverityNumber(*sd.SeverityNumber())
		}

		logRecord := &logspb.LogRecord{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(sd.ObservedTimestamp().UnixNano()),
			TraceId:              traceIDBytes,       // provide the associated trace ID if available
			SpanId:               spanIDBytes,        // provide the associated span ID if available
			Flags:                uint32(traceFlags), // provide the associated trace flags
			Body:                 body,               // provide the associated log body if available
			Attributes:           kv,                 // provide additional log attributes if available
			SeverityText:         st,
			SeverityNumber:       sn,
		}

		var is *commonpb.InstrumentationScope
		var schemaURL = ""
		if sd.InstrumentationScope() != nil {
			is = &commonpb.InstrumentationScope{
				Name:    sd.InstrumentationScope().Name,
				Version: sd.InstrumentationScope().Version,
			}
			schemaURL = sd.InstrumentationScope().SchemaURL
		}

		// Create a log resource
		resourceLog := &logspb.ResourceLogs{
			Resource: &resourcepb.Resource{
				Attributes: KeyValues(sd.Resource().Attributes()),
			},
			// provide a resource description if available
			ScopeLogs: []*logspb.ScopeLogs{
				{
					Scope:      is,
					SchemaUrl:  schemaURL,
					LogRecords: []*logspb.LogRecord{logRecord},
				},
			},
		}

		resourceLogs = append(resourceLogs, resourceLog)
	}

	return resourceLogs
}
