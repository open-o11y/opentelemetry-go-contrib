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
)

func (t traceID) convertTraceIDToHexString() string {
	return hex.EncodeToString(t[:])
}

func (s spanID) cconvertSpanIDToHexString() string {
	return hex.EncodeToString(s[:])
}

func TestAwsXRayTraceIdIsValidLength(t *testing.T) {
	idg := awsXRayIDGenerator()
	traceIDHex := idg.NewTraceID().convertTraceIDToHexString()
	traceIDLength := len(traceIDHex)

	if traceIDLength != 32 {
		t.Errorf("TraceID has incorrect length. Got length of %d, expected 32", traceIDLength)
	}
}

func TestAwsXRayTraceIdIsUnique(t *testing.T) {
	idg := awsXRayIDGenerator()
	traceID1 := idg.NewTraceID().convertTraceIDToHexString()
	traceID2 := idg.NewTraceID().convertTraceIDToHexString()

	if traceID1 == traceID2 {
		t.Errorf("TraceID should be unique. Got TraceID1 = %s and TraceID2 = %s", traceID1, traceID2)
	}
}

func TestAwsXRayTraceIdTimeStampInBounds(t *testing.T) {

	idg := awsXRayIDGenerator()

	previousTime := time.Now().Unix()

	traceIDHex := idg.NewTraceID().convertTraceIDToHexString()
	currentTime, err := strconv.ParseInt(traceIDHex[0:8], 16, 64)

	nextTime := time.Now().Unix()

	if err != nil {
		t.Error(err)
	}

	inLowerBound := previousTime <= currentTime
	inUpperBound := currentTime <= nextTime

	if !inLowerBound || !inUpperBound {
		t.Errorf("TraceID is generated incorrectly with the wrong timestamp. Got epoch time %d, expected epoch time should be between %d and %d", currentTime, previousTime, currentTime)
	}
}

func TestAwsXRayTraceIdIsNotNil(t *testing.T) {
	var nilTraceID traceID
	idg := awsXRayIDGenerator()
	traceID := idg.NewTraceID()
	isNil := bytes.Equal(traceID[:], nilTraceID[:])

	if isNil == true {
		t.Error("TraceID cannot be Nil.")
	}
}

func TestAwsXRaySpanIdIsValidLength(t *testing.T) {
	idg := awsXRayIDGenerator()
	spanIDHex := idg.NewSpanID().cconvertSpanIDToHexString()
	spanIDLength := len(spanIDHex)

	if spanIDLength != 16 {
		t.Errorf("SpanID has incorrect length. Got length of %d, expected 16", spanIDLength)
	}
}

func TestAwsXRaySpanIdIsUnique(t *testing.T) {
	idg := awsXRayIDGenerator()
	spanID1 := idg.NewSpanID().cconvertSpanIDToHexString()
	spanID2 := idg.NewSpanID().cconvertSpanIDToHexString()

	if spanID1 == spanID2 {
		t.Errorf("SpanID should be unique. Got spanID1 = %s and spanID2 = %s", spanID1, spanID2)
	}
}

func TestAwsXRaySpanIdIsNotNil(t *testing.T) {
	var nilSpanID spanID
	idg := awsXRayIDGenerator()
	spanID := idg.NewSpanID()
	isNil := bytes.Equal(spanID[:], nilSpanID[:])

	if isNil == true {
		t.Error("SpanID cannot be Nil.")
	}
}
