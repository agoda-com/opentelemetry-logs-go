package logs

import (
	"context"
	"go.opentelemetry.io/otel"
	"log"
	"sync"
)

type simpleLogRecordProcessor struct {
	exporterMu sync.Mutex
	stopOnce   sync.Once
	exporter   LogRecordExporter
}

func (lrp *simpleLogRecordProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (lrp *simpleLogRecordProcessor) ForceFlush(ctx context.Context) error {
	return nil
}

var _ LogRecordProcessor = (*simpleLogRecordProcessor)(nil)

// NewSimpleLogRecordProcessor returns a new LogRecordProcessor that will synchronously
// send completed logs to the exporter immediately.
//
// This LogRecordProcessor is not recommended for production use. The synchronous
// nature of this LogRecordProcessor make it good for testing, debugging, or
// showing examples of other feature, but it will be slow and have a high
// computation resource usage overhead. The BatchLogsProcessor is recommended
// for production use instead.
func NewSimpleLogRecordProcessor(exporter LogRecordExporter) LogRecordProcessor {
	slp := &simpleLogRecordProcessor{
		exporter: exporter,
	}
	log.Printf("SimpleLogsProcessor is not recommended for production use, consider using BatchSpanProcessor instead.")

	return slp
}

// OnEmit Process immediately emits a LogRecord
func (lrp *simpleLogRecordProcessor) OnEmit(rol ReadableLogRecord) {
	lrp.exporterMu.Lock()
	defer lrp.exporterMu.Unlock()

	if err := lrp.exporter.Export(context.Background(), []ReadableLogRecord{rol}); err != nil {
		otel.Handle(err)
	}
}

// MarshalLog is the marshaling function used by the logging system to represent this LogRecord Processor.
func (lrp *simpleLogRecordProcessor) MarshalLog() interface{} {
	return struct {
		Type     string
		Exporter LogRecordExporter
	}{
		Type:     "SimpleLogRecordProcessor",
		Exporter: lrp.exporter,
	}
}
