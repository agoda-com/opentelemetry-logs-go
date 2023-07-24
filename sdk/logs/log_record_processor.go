package logs

import (
	"context"
	"sync"
)

// LogRecordProcessor is an interface which allows hooks for LogRecord emitting.
// see https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor
type LogRecordProcessor interface {

	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// OnEmit is called when logs sent. It is called synchronously and
	// hence not block.
	OnEmit(rol ReadableLogRecord)
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the processor should be done in this call.
	//
	// Calls to Process, or ForceFlush after this has been called
	// should be ignored.
	//
	// All timeouts and cancellations contained in ctx must be honored, this
	// should not block indefinitely.
	Shutdown(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// ForceFlush exports all ended logs to the configured Exporter that have not yet
	// been exported.  It should only be called when absolutely necessary, such as when
	// using a FaaS provider that may suspend the process after an invocation, but before
	// the Processor can export the completed spans.
	ForceFlush(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

type logRecordProcessorState struct {
	lp    LogRecordProcessor
	state sync.Once
}

func newLogsProcessorState(lp LogRecordProcessor) *logRecordProcessorState {
	return &logRecordProcessorState{lp: lp}
}

type logRecordProcessorStates []*logRecordProcessorState
