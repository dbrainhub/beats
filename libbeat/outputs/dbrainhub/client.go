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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/publisher"
)

// Client is a dbrainhub client.
type dbRainhubClient struct {
	client   *http.Client
	endpoint string
	observer outputs.Observer
	timeout  time.Duration
	log      *logp.Logger
	dbIp     string
	dbPort   int
}

type DbRainhubResponse struct {
	FailedEvents []int32 `json:"failed_events,omitempty"`
}

func (hc *dbRainhubClient) Publish(ctx context.Context, batch publisher.Batch) error {
	events := batch.Events()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(map[string]interface{}{
		"events":  events,
		"db_ip":   hc.dbIp,
		"db_port": hc.dbPort,
	})
	if err != nil {
		hc.log.Errorf("Failed to encode events: %v", events)
		return err
	}

	req, err := http.NewRequest("POST", hc.endpoint, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := hc.client.Do(req)
	if err != nil {
		// error type
		// TODO: limiter

		// TODO: not classified

		return err
	} else {
		status := resp.StatusCode
		if status == 200 {
			// get and check failed events
			var dbRainhubRes DbRainhubResponse
			err = json.NewDecoder(resp.Body).Decode(&dbRainhubRes)
			if err != nil {
				return err
			}
			failedEvents := dbRainhubRes.FailedEvents
			failedNo := len(failedEvents)
			if failedNo == 0 {
				// all success
				hc.observer.Acked(len(events))
				batch.ACK()
			} else {
				// retry the failed events
				hc.observer.Failed(failedNo)
				hc.observer.Acked(len(events) - failedNo)
				retryEvents := make([]publisher.Event, failedNo)
				for _, i := range failedEvents {
					retryEvents = append(retryEvents, events[i])
				}
				batch.RetryEvents(retryEvents)
			}
		} else {
			// failed and retry
			batch.RetryEvents(events)
		}
	}

	if resp != nil {
		defer resp.Body.Close()
		_, _ = ioutil.ReadAll(resp.Body)
	}
	return nil
}

func (hc *dbRainhubClient) Connect() error {
	hc.client = &http.Client{
		Timeout: hc.timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Timeout: hc.timeout}).DialContext,
		},
	}
	return nil
}

func (hc *dbRainhubClient) Close() error {
	hc.client = nil
	return nil
}

func (hc *dbRainhubClient) String() string {
	return "dbrainhub-" + hc.endpoint
}
