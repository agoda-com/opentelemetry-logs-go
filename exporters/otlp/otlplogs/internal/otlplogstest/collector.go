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

package otlplogstest

import (
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	"sort"

	collectorlogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

// LogsCollector mocks a collector for the end-to-end testing.
type LogsCollector interface {
	Stop() error
	GetResourceLogs() []*logspb.ResourceLogs
}

// LogsStorage stores the Logs. Mock collectors can use it to
// store logs they have received.
type LogsStorage struct {
	rlm       map[string]*logspb.ResourceLogs
	logsCount int
}

// NewLogsStorage creates a new logs storage.
func NewLogsStorage() LogsStorage {
	return LogsStorage{
		rlm: make(map[string]*logspb.ResourceLogs),
	}
}

// AddLogs adds logs to the logs storage.
func (s *LogsStorage) AddLogs(request *collectorlogspb.ExportLogsServiceRequest) {
	for _, rs := range request.GetResourceLogs() {
		rstr := resourceString(rs.Resource)
		if existingRs, ok := s.rlm[rstr]; !ok {
			s.rlm[rstr] = rs
			// TODO (rghetia): Add support for library Info.
			if len(rs.ScopeLogs) == 0 {
				rs.ScopeLogs = []*logspb.ScopeLogs{
					{
						LogRecords: []*logspb.LogRecord{},
					},
				}
			}
			s.logsCount += len(rs.ScopeLogs[0].LogRecords)
		} else {
			if len(rs.ScopeLogs) > 0 {
				newLogs := rs.ScopeLogs[0].GetLogRecords()
				existingRs.ScopeLogs[0].LogRecords = append(existingRs.ScopeLogs[0].LogRecords, newLogs...)
				s.logsCount += len(newLogs)
			}
		}
	}
}

// GetLogRecords returns the stored logs.
func (s *LogsStorage) GetLogRecords() []*logspb.LogRecord {
	logs := make([]*logspb.LogRecord, 0, s.logsCount)
	for _, rs := range s.rlm {
		logs = append(logs, rs.ScopeLogs[0].LogRecords...)
	}
	return logs
}

// GetResourceLogs returns the stored resource logs.
func (s *LogsStorage) GetResourceLogs() []*logspb.ResourceLogs {
	rls := make([]*logspb.ResourceLogs, 0, len(s.rlm))
	for _, rs := range s.rlm {
		rls = append(rls, rs)
	}
	return rls
}

func resourceString(res *resourcepb.Resource) string {
	sAttrs := sortedAttributes(res.GetAttributes())
	rstr := ""
	for _, attr := range sAttrs {
		rstr = rstr + attr.String()
	}
	return rstr
}

func sortedAttributes(attrs []*commonpb.KeyValue) []*commonpb.KeyValue {
	sort.Slice(attrs[:], func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
}
