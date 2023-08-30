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
	logssdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/agoda-com/opentelemetry-logs-go/sdk/logs/logstest"
	"github.com/stretchr/testify/assert"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	"testing"
	"time"
)

func TestNilLogRecord(t *testing.T) {
	assert.Nil(t, Logs(nil))
}

func TestEmptyLogRecord(t *testing.T) {
	assert.Nil(t, Logs([]logssdk.ReadableLogRecord{}))
}

func TestLogRecord(t *testing.T) {
	//attrs := []attribute.KeyValue{attribute.Int("one", 1), attribute.Int("two", 2)}
	//eventTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)

	logTime := time.Unix(1589932800, 0000)
	body := "message"

	lr := logRecord(logstest.LogRecordStub{
		Timestamp:         &logTime,
		ObservedTimestamp: logTime,
		Body:              &body,
	}.Snapshot())

	logTimestamp := uint64(1589932800 * 1e9)

	bodyValue := &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "message"}}
	assert.Equal(t, &logspb.LogRecord{
		Attributes:           nil,
		Body:                 bodyValue,
		TimeUnixNano:         logTimestamp,
		ObservedTimeUnixNano: logTimestamp,
	}, lr)
}
