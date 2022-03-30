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

package rainhub

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

// Client is a rainhub client.
type rainhubClient struct {
	client				*http.Client
	endpoint 			string
	observer            outputs.Observer
	timeout				time.Duration
	log 				*logp.Logger
}

func (hc *rainhubClient) Publish(ctx context.Context, batch publisher.Batch) error {
	events := batch.Events()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(events)
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
		batch.RetryEvents(events)
	} else {
		hc.observer.Acked(len(events))
		batch.ACK()
	}

	if resp != nil {
		defer resp.Body.Close()
		_, _ = ioutil.ReadAll(resp.Body)
	}
	return nil
}

func (hc *rainhubClient) Connect() error {
	hc.client = &http.Client{
		Timeout: hc.timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Timeout: hc.timeout}).DialContext,
		},
	}
	return nil
}

func (hc *rainhubClient) Close() error {
	hc.client = nil
	return nil
}

func (hc *rainhubClient) String() string {
	return "rainhub-" + hc.endpoint
}
