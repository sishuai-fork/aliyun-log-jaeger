// Copyright (c) 2018 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanstore

import (
	"testing"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSpanReader_logsToStringArray(t *testing.T) {
	logs := make([]map[string]string, 3)
	logs[0] = map[string]string{
		traceIDField:       "0",
		operationNameField: "op0",
	}
	logs[1] = map[string]string{
		traceIDField:       "1",
		operationNameField: "op1",
	}
	logs[2] = map[string]string{
		traceIDField:       "2",
		operationNameField: "op2",
	}
	actual, err := logsToStringArray(logs, operationNameField)
	require.NoError(t, err)
	assert.EqualValues(t, []string{"op0", "op1", "op2"}, actual)
}

func TestSpanReader_buildFindTraceIDsQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `where "process.serviceName" = 's' and operationName = 'o' and 1000000000 <= duration and duration <= 2000000000 and "tags.http.status_code" = '200'`
	traceQuery := &spanstore.TraceQueryParameters{
		DurationMin:   time.Second,
		DurationMax:   time.Second * 2,
		StartTimeMin:  time.Time{},
		StartTimeMax:  time.Time{}.Add(time.Second),
		ServiceName:   "s",
		OperationName: "o",
		Tags: map[string]string{
			"http.status_code": "200",
		},
	}
	actualQuery := r.buildFindTraceIDsQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildFindTraceIDsQuery_emptyQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := ""
	traceQuery := &spanstore.TraceQueryParameters{}
	actualQuery := r.buildFindTraceIDsQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildFindTraceIDsQuery_singleCondition(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `where "process.serviceName" = 'svc1'`
	traceQuery := &spanstore.TraceQueryParameters{
		ServiceName: "svc1",
	}
	actualQuery := r.buildFindTraceIDsQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildServiceNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `"process.serviceName" = 'svc1'`
	serviceNameQuery := r.buildServiceNameQuery("svc1")
	assert.Equal(t, expectedStr, serviceNameQuery)
}

func TestSpanReader_buildOperationNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := "operationName = 'op1'"
	operationNameQuery := r.buildOperationNameQuery("op1")
	assert.Equal(t, expectedStr, operationNameQuery)
}

func TestSpanReader_buildDurationQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := "1000000000 <= duration and duration <= 2000000000"
	durationMin := time.Second
	durationMax := time.Second * 2
	durationQuery := r.buildDurationQuery(durationMin, durationMax)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = "12000000 <= duration"
	durationMin = time.Millisecond * 12
	durationQuery = r.buildDurationQuery(durationMin, 0)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = "duration <= 18000000000000"
	durationMax = time.Hour * 5
	durationQuery = r.buildDurationQuery(0, durationMax)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = ""
	durationQuery = r.buildDurationQuery(0, 0)
	assert.Equal(t, expectedStr, durationQuery)
}

func TestSpanReader_buildTagQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `"tags.http.method" = 'POST'`
	operationNameQuery := r.buildTagQuery("http.method", "POST")
	assert.Equal(t, expectedStr, operationNameQuery)
}