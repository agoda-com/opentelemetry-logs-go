package logs

import (
	"context"
	"sync"
)

type testExporter struct {
	mu   sync.RWMutex
	idx  map[string]int
	logs []*ReadableLogRecord
}

func NewTestExporter() *testExporter {
	return &testExporter{idx: make(map[string]int)}
}

func (te *testExporter) Export(ctx context.Context, logs []ReadableLogRecord) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	i := len(te.logs)
	for _, s := range logs {
		te.logs = append(te.logs, &s)
		i++
	}
	return nil
}

func (te *testExporter) Shutdown(ctx context.Context) error {
	return nil
}
