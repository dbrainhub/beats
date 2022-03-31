// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package dbrainhub

import (
	"context"
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/outest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientWithHeaders(t *testing.T) {
	requestCount := 0
	// start a mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"items":[{"index":{}},{"index":{}},{"index":{}}]}`
		fmt.Fprintln(w, response)
		requestCount++
	}))
	defer ts.Close()

	client := &dbRainhubClient{
		observer:	outputs.NewNilObserver(),
		endpoint: 	ts.URL,
		timeout:  	time.Duration(2000000000),
		log: 	  	logp.NewLogger("dbrainhub"),
	}
	client.Connect()
	defer  client.Close()

	// publish
	event := beat.Event{Fields: common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "dbrainhub",
		"message":    "Test message from dbrainhub",
	}}
	batch := outest.NewBatch(event)
	err := client.Publish(context.Background(), batch)
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount)
}
