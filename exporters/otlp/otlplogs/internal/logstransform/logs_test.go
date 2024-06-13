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
	"math"
	"testing"
	"time"
)

func TestNilLogRecord(t *testing.T) {
	assert.Nil(t, Logs(nil))
}

func TestEmptyLogRecord(t *testing.T) {
	assert.Nil(t, Logs([]logssdk.ReadableLogRecord{}))
}

type logStruct struct {
	StringField     string
	IntField        int
	FloatField      float64 `json:"float_field"`
	UintField       uint
	BoolField       bool
	TimeField       time.Time
	BytesField      []byte
	BytesArrayField [8]byte
	SliceField      []logStruct
	ArrayField      [1]*logStruct
	MapField        map[string]any
	StructField     *logStruct
	NilField        *logStruct
}

func TestLogRecord(t *testing.T) {
	//attrs := []attribute.KeyValue{attribute.Int("one", 1), attribute.Int("two", 2)}
	//eventTime := time.Date(2020, 5, 20, 0, 0, 0, 0, time.UTC)

	logTime := time.Unix(1589932800, 0000)
	refTime := time.Now()
	body := logStruct{
		StringField:     "hello world",
		IntField:        123,
		FloatField:      3.14,
		UintField:       math.MaxUint64,
		BoolField:       true,
		TimeField:       refTime,
		BytesField:      []byte{1, 3, 5},
		BytesArrayField: [8]byte{1, 3, 5},
		SliceField: []logStruct{
			{
				StringField: "sliced_field_nested",
			},
		},
		ArrayField: [1]*logStruct{},
		MapField: map[string]any{
			"first_key": "first_value",
			"nested_key": map[string]any{
				"second_key": "second_value",
			},
		},
		StructField: &logStruct{
			StringField: "hello world2",
		},
		NilField: nil,
	}

	lr := logRecord(logstest.LogRecordStub{
		Timestamp:         &logTime,
		ObservedTimestamp: logTime,
		Body:              &body,
	}.Snapshot())

	logTimestamp := uint64(1589932800 * 1e9)

	bodyValue := &commonpb.AnyValue{
		Value: &commonpb.AnyValue_KvlistValue{
			KvlistValue: &commonpb.KeyValueList{
				Values: []*commonpb.KeyValue{
					{
						Key: "StringField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_StringValue{
								StringValue: body.StringField,
							},
						},
					},
					{
						Key: "IntField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_IntValue{
								IntValue: int64(body.IntField),
							},
						},
					},
					{
						Key: "float_field",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_DoubleValue{
								DoubleValue: body.FloatField,
							},
						},
					},
					{
						Key: "UintField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_IntValue{
								IntValue: math.MaxInt64,
							},
						},
					},
					{
						Key: "BoolField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_BoolValue{
								BoolValue: body.BoolField,
							},
						},
					},
					{
						Key: "TimeField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_StringValue{
								StringValue: body.TimeField.Format(time.RFC3339Nano),
							},
						},
					},
					{
						Key: "BytesField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_BytesValue{
								BytesValue: body.BytesField,
							},
						},
					},
					{
						Key: "BytesArrayField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_BytesValue{
								BytesValue: body.BytesArrayField[:],
							},
						},
					},
					{
						Key: "SliceField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_ArrayValue{
								ArrayValue: &commonpb.ArrayValue{
									Values: []*commonpb.AnyValue{
										{
											Value: &commonpb.AnyValue_KvlistValue{
												KvlistValue: &commonpb.KeyValueList{
													Values: []*commonpb.KeyValue{
														{
															Key: "StringField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_StringValue{
																	StringValue: body.SliceField[0].StringField,
																},
															},
														},
														{
															Key: "IntField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_IntValue{
																	IntValue: int64(body.SliceField[0].IntField),
																},
															},
														},
														{
															Key: "float_field",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_DoubleValue{
																	DoubleValue: body.SliceField[0].FloatField,
																},
															},
														},
														{
															Key: "UintField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_IntValue{
																	IntValue: int64(body.SliceField[0].UintField),
																},
															},
														},
														{
															Key: "BoolField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_BoolValue{
																	BoolValue: body.SliceField[0].BoolField,
																},
															},
														},
														{
															Key: "BytesArrayField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_BytesValue{
																	BytesValue: body.SliceField[0].BytesArrayField[:],
																},
															},
														},
														{
															Key: "ArrayField",
															Value: &commonpb.AnyValue{
																Value: &commonpb.AnyValue_ArrayValue{
																	ArrayValue: &commonpb.ArrayValue{
																		Values: []*commonpb.AnyValue{
																			nil,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Key: "ArrayField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_ArrayValue{
								ArrayValue: &commonpb.ArrayValue{
									Values: []*commonpb.AnyValue{
										nil,
									},
								},
							},
						},
					},
					{
						Key: "MapField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_KvlistValue{
								KvlistValue: &commonpb.KeyValueList{
									Values: []*commonpb.KeyValue{
										{
											Key: "first_key",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: "first_value",
												},
											},
										},
										{
											Key: "nested_key",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_KvlistValue{
													KvlistValue: &commonpb.KeyValueList{
														Values: []*commonpb.KeyValue{
															{
																Key: "second_key",
																Value: &commonpb.AnyValue{
																	Value: &commonpb.AnyValue_StringValue{
																		StringValue: "second_value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Key: "StructField",
						Value: &commonpb.AnyValue{
							Value: &commonpb.AnyValue_KvlistValue{
								KvlistValue: &commonpb.KeyValueList{
									Values: []*commonpb.KeyValue{
										{
											Key: "StringField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_StringValue{
													StringValue: body.StructField.StringField,
												},
											},
										},
										{
											Key: "IntField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_IntValue{
													IntValue: int64(body.StructField.IntField),
												},
											},
										},
										{
											Key: "float_field",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_DoubleValue{
													DoubleValue: body.StructField.FloatField,
												},
											},
										},
										{
											Key: "UintField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_IntValue{
													IntValue: int64(body.StructField.UintField),
												},
											},
										},
										{
											Key: "BoolField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BoolValue{
													BoolValue: body.StructField.BoolField,
												},
											},
										},
										{
											Key: "BytesArrayField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_BytesValue{
													BytesValue: body.StructField.BytesArrayField[:],
												},
											},
										},
										{
											Key: "ArrayField",
											Value: &commonpb.AnyValue{
												Value: &commonpb.AnyValue_ArrayValue{
													ArrayValue: &commonpb.ArrayValue{
														Values: []*commonpb.AnyValue{
															nil,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, &logspb.LogRecord{
		Attributes:           nil,
		Body:                 bodyValue,
		TimeUnixNano:         logTimestamp,
		ObservedTimeUnixNano: logTimestamp,
	}, lr)
}
