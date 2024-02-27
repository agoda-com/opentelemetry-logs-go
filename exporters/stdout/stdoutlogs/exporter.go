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
	"encoding/json"
	"sync"

	sdk "github.com/kudarap/opentelemetry-logs-go/sdk/logs"
)

var _ sdk.LogRecordExporter = &Exporter{}

// NewExporter creates an Exporter with the passed options.
func NewExporter(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	return &Exporter{
		encoder: enc,
	}, nil
}

// Exporter is an implementation of logs.LogRecordSyncer that writes spans to stdout.
type Exporter struct {
	encoder   *json.Encoder
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

	logRecords := logRecordsFromReadableLogRecords(logs)

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	for _, lr := range logRecords {
		// Encode span stubs, one by one
		if err := e.encoder.Encode(lr); err != nil {
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
