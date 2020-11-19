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

package aws // import "go.opentelemetry.io/contrib/idgenerator/aws"

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type spanID [8]byte
type traceID [16]byte

// IDGenerator is an interface for generating new TraceIDs and SpanIDs
type IDGenerator interface {
	NewTraceID() traceID
	NewSpanID() spanID
}

type xRayIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// XRayIDGenerator returns an idGenerator used for sending traces to AWS X-Ray
func XRayIDGenerator() IDGenerator {
	gen := &xRayIDGenerator{}
	var rngSeed int64
	err := binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	if err != nil {
		panic(err)
	}
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *xRayIDGenerator) NewSpanID() spanID {
	gen.Lock()
	defer gen.Unlock()
	sid := spanID{}
	gen.randSource.Read(sid[:])
	return sid
}

// NewTraceID returns a non-zero trace ID based on AWS X-Ray TraceID format.
// (https://docs.aws.amazon.com/xray/latest/devguide/xray-api-sendingdata.html#xray-api-traceids)
func (gen *xRayIDGenerator) NewTraceID() traceID {
	gen.Lock()
	defer gen.Unlock()

	tid := traceID{}
	currentTime := getCurrentTimeHex()
	copy(tid[:4], currentTime)
	gen.randSource.Read(tid[4:])

	return tid
}

func getCurrentTimeHex() []uint8 {
	currentTime := time.Now().Unix()
	currentTimeHex, err := hex.DecodeString(strconv.FormatInt(currentTime, 16))
	if err != nil {
		panic(err)
	}
	return currentTimeHex
}
