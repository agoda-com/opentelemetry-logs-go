package logstransform

import (
	sdk "github.com/agoda-com/otellogs-go/otel/sdk/logs"
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

		logRecord := &logspb.LogRecord{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(sd.ObservedTimestamp().UnixNano()),
			TraceId:              traceIDBytes,                // provide the associated trace ID if available
			SpanId:               spanIDBytes,                 // provide the associated span ID if available
			Flags:                uint32(traceFlags),          // provide the associated trace flags
			Body:                 body,                        // provide the associated log body if available
			Attributes:           KeyValues(*sd.Attributes()), // provide additional log attributes if available
			SeverityText:         *sd.SeverityText(),
			SeverityNumber:       logspb.SeverityNumber(*sd.SeverityNumber()),
		}

		// Create a log resource
		resourceLog := &logspb.ResourceLogs{
			Resource: &resourcepb.Resource{
				Attributes: KeyValues(sd.Resource().Attributes()),
			},
			// provide a resource description if available
			ScopeLogs: []*logspb.ScopeLogs{
				{
					Scope: &commonpb.InstrumentationScope{
						Name:    sd.InstrumentationScope().Name,
						Version: sd.InstrumentationScope().Version,
					},
					SchemaUrl:  sd.InstrumentationScope().SchemaURL,
					LogRecords: []*logspb.LogRecord{logRecord},
				},
			},
		}

		resourceLogs = append(resourceLogs, resourceLog)
	}

	return resourceLogs
}
