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

package stdoutlogs

import (
	"context"
	"fmt"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var _ sdk.LogRecordExporter = &Exporter{}

// NewExporter creates an Exporter with the passed options.
func NewExporter(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		writer: cfg.Writer,
	}, nil
}

// Exporter is an implementation of logs.LogRecordSyncer that writes spans to stdout.
type Exporter struct {
	writer    io.Writer
	encoderMu sync.Mutex

	stoppedMu sync.RWMutex
	stopped   bool
}

// Export writes logs in json format to stdout.
func (e *Exporter) Export(ctx context.Context, logs []sdk.ReadableLogRecord) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(logs) == 0 {
		return nil
	}

	wr := os.Stdout

	logRecords := logRecordsFromReadableLogRecords(logs)

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	for _, lr := range logRecords {

		var logMessageBuilder strings.Builder

		logMessageBuilder.WriteString(lr.ObservedTimestamp.Format(time.RFC3339))
		logMessageBuilder.WriteString(" ")
		logMessageBuilder.WriteString(lr.getSeverityText())
		logMessageBuilder.WriteString(" ")
		if lr.Body != nil {
			logMessageBuilder.WriteString(*lr.Body)
			logMessageBuilder.WriteString(" ")
		}

		if lr.TraceId != nil || lr.SpanId != nil {
			logMessageBuilder.WriteString(": ")
			if lr.TraceId != nil {
				traceId := *lr.TraceId
				logMessageBuilder.WriteString("traceId=")
				logMessageBuilder.WriteString(traceId.String())
				logMessageBuilder.WriteString(" ")
			}
			if lr.SpanId != nil {
				spanId := *lr.SpanId
				logMessageBuilder.WriteString("spanId=")
				logMessageBuilder.WriteString(spanId.String())
				logMessageBuilder.WriteString(" ")
			}
		}

		if lr.InstrumentationScope != nil {
			logMessageBuilder.WriteString("[scopeInfo: ")
			scope := *lr.InstrumentationScope
			logMessageBuilder.WriteString(scope.Name)
			if scope.Version != "" {
				logMessageBuilder.WriteString(":")
				logMessageBuilder.WriteString(scope.Version)
			}
			logMessageBuilder.WriteString("] ")
		}

		attributes := lr.Resource.Attributes()
		if lr.Attributes != nil {
			attributes = append(attributes, *lr.Attributes...)
		}

		if len(attributes) > 0 {
			logMessageBuilder.WriteString("{")
			for i, a := range attributes {
				logMessageBuilder.WriteString(string(a.Key))
				logMessageBuilder.WriteString("=")
				logMessageBuilder.WriteString(a.Value.AsString())
				if i < len(attributes)-1 {
					logMessageBuilder.WriteString(", ")
				}
			}
			logMessageBuilder.WriteString("}")
		}

		// Print logRecords, one by one
		_, err := fmt.Fprintf(wr, "%s\n", logMessageBuilder.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown is called to stop the exporter, it performs no action.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (e *Exporter) MarshalLog() interface{} {
	return struct {
		Type           string
		WithTimestamps bool
	}{
		Type: "stdout",
	}
}
