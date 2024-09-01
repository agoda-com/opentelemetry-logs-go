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
	"fmt"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

// stdOutLogRecord is a stand-in for a LogRecord.
type stdOutLogRecord struct {
	Timestamp            *time.Time
	ObservedTimestamp    time.Time
	TraceId              *trace.TraceID
	SpanId               *trace.SpanID
	TraceFlags           *trace.TraceFlags
	SeverityText         *string
	SeverityNumber       *logs.SeverityNumber
	Body                 *string
	Resource             *resource.Resource
	InstrumentationScope *instrumentation.Scope
	Attributes           *[]attribute.KeyValue
}

func (lr stdOutLogRecord) getSeverityText() string {
	if lr.SeverityNumber == nil {
		return "UNSPECIFIED"
	} else {
		sn := int32(*lr.SeverityNumber)
		if sn >= 1 && sn <= 4 {
			return "TRACE"
		}
		if sn >= 5 && sn <= 8 {
			return "DEBUG"
		}
		if sn >= 9 && sn <= 12 {
			return "INFO"
		}
		if sn >= 13 && sn <= 16 {
			return "WARN"
		}
		if sn >= 17 && sn <= 20 {
			return "ERROR"
		}
		if sn >= 21 && sn <= 24 {
			return "FATAL"
		}
		return "UNSPECIFIED"
	}
}

func logRecordsFromReadableLogRecords(logRecords []sdk.ReadableLogRecord) []stdOutLogRecord {

	var result []stdOutLogRecord

	for _, lr := range logRecords {
		logRecord := stdOutLogRecord{
			Timestamp:            lr.Timestamp(),
			ObservedTimestamp:    lr.ObservedTimestamp(),
			TraceId:              lr.TraceId(),
			SpanId:               lr.SpanId(),
			TraceFlags:           lr.TraceFlags(),
			SeverityText:         lr.SeverityText(),
			SeverityNumber:       lr.SeverityNumber(),
			Body:                 convertBodyToString(lr.Body()),
			Resource:             lr.Resource(),
			InstrumentationScope: lr.InstrumentationScope(),
			Attributes:           lr.Attributes(),
		}
		result = append(result, logRecord)
	}
	return result
}

func convertBodyToString(body any) *string {
	typ := reflect.TypeOf(body)
	val := reflect.ValueOf(body)
	if valueIsNil(typ, val) {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if !val.CanInterface() {
		return nil
	}
	var str string
	switch {
	case typ.ConvertibleTo(reflect.TypeOf(time.Time{})):
		valTime := val.Convert(reflect.TypeOf(time.Time{})).Interface().(time.Time)
		str = valTime.Format(time.RFC3339Nano)
	case typ.Kind() == reflect.Struct:
		str = fmt.Sprintf("%+v", val.Interface())
	default:
		str = fmt.Sprintf("%v", val.Interface())
	}
	return &str
}

func valueIsNil(typ reflect.Type, val reflect.Value) bool {
	if typ == nil {
		return true
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}
