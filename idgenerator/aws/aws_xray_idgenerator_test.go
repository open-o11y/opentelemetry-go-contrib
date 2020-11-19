// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (t traceID) convertTraceIDToHexString() string {
	return hex.EncodeToString(t[:])
}

func (s spanID) convertSpanIDToHexString() string {
	return hex.EncodeToString(s[:])
}

func TestAwsXRayTraceIdIsValidLength(t *testing.T) {
	idg := XRayIDGenerator()
	traceIDHex := idg.NewTraceID().convertTraceIDToHexString()
	traceIDLength := len(traceIDHex)
	expectedTraceIDLength := 32

	assert.Equal(t, traceIDLength, expectedTraceIDLength, "TraceID has incorrect length.")
}

func TestAwsXRayTraceIdIsUnique(t *testing.T) {
	idg := XRayIDGenerator()
	traceID1 := idg.NewTraceID().convertTraceIDToHexString()
	traceID2 := idg.NewTraceID().convertTraceIDToHexString()

	assert.NotEqual(t, traceID1, traceID2, "TraceID should be unique")
}

func TestAwsXRayTraceIdTimeStampInBounds(t *testing.T) {

	idg := XRayIDGenerator()

	previousTime := time.Now().Unix()

	traceIDHex := idg.NewTraceID().convertTraceIDToHexString()
	currentTime, err := strconv.ParseInt(traceIDHex[0:8], 16, 64)

	nextTime := time.Now().Unix()

	if err != nil {
		t.Error(err)
	}

	assert.LessOrEqual(t, previousTime, currentTime, "TraceID is generated incorrectly with the wrong timestamp.")
	assert.LessOrEqual(t, currentTime, nextTime, "TraceID is generated incorrectly with the wrong timestamp.")
}

func TestAwsXRayTraceIdIsNotNil(t *testing.T) {
	var nilTraceID traceID
	idg := XRayIDGenerator()
	traceID := idg.NewTraceID()

	assert.False(t, bytes.Equal(traceID[:], nilTraceID[:]), "TraceID cannot be Nil.")
}

func TestAwsXRaySpanIdIsValidLength(t *testing.T) {
	idg := XRayIDGenerator()
	spanIDHex := idg.NewSpanID().convertSpanIDToHexString()
	spanIDLength := len(spanIDHex)
	expectedSpanIDLength := 16

	assert.Equal(t, spanIDLength, expectedSpanIDLength, "SpanID has incorrect length")
}

func TestAwsXRaySpanIdIsUnique(t *testing.T) {
	idg := XRayIDGenerator()
	spanID1 := idg.NewSpanID().convertSpanIDToHexString()
	spanID2 := idg.NewSpanID().convertSpanIDToHexString()

	assert.NotEqual(t, spanID1, spanID2, "SpanID should be unique")
}

func TestAwsXRaySpanIdIsNotNil(t *testing.T) {
	var nilSpanID spanID
	idg := XRayIDGenerator()
	spanID := idg.NewSpanID()

	assert.False(t, bytes.Equal(spanID[:], nilSpanID[:]), "SpanID cannot be Nil.")
}
