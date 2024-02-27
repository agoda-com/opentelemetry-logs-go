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

package otel

import (
	"github.com/kudarap/opentelemetry-logs-go/logs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testLoggerProvider struct{}

var _ logs.LoggerProvider = &testLoggerProvider{}

func (*testLoggerProvider) Logger(_ string, _ ...logs.LoggerOption) logs.Logger {
	return logs.NewNoopLoggerProvider().Logger("")
}

func TestMultipleGlobalLoggerProvider(t *testing.T) {
	p1 := testLoggerProvider{}
	p2 := logs.NewNoopLoggerProvider()
	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	got := GetLoggerProvider()
	assert.Equal(t, p2, got)
}
