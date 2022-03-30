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
	"testing"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/stretchr/testify/assert"
)

func TestMakeRH(t *testing.T) {
	config := `
hosts: ["127.0.0.1:8080"]
batch_size: 10
retry_limit: 5
timeout: 2
`
	cfg := common.MustNewConfigFrom(config)
	beatInfo := beat.Info{Beat: "libbeat", Version: "1.2.3"}
	_, err := makeRH(nil, beatInfo, outputs.NewNilObserver(), cfg)
	assert.NoError(t, err)
}

func TestMakeRHError(t *testing.T) {
	config := `
hosts:
batch_size: 10
retry_limit: 5
timeout: 2
`
	cfg := common.MustNewConfigFrom(config)
	beatInfo := beat.Info{Beat: "libbeat", Version: "1.2.3"}
	_, err := makeRH(nil, beatInfo, outputs.NewNilObserver(), cfg)
	assert.Error(t, err)
}

