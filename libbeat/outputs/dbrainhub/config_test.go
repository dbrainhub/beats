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
	"testing"
	"time"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	config := `
hosts: ["127.0.0.1:8080"]
batch_size: 10
retry_limit: 5
timeout: 2
`
	c := common.MustNewConfigFrom(config)
	rhConfig, err := readConfig(c)
	if err != nil {
		t.Fatalf("Can't create test configuration from valid input")
	}

	assert.Equal(t, []string([]string{"127.0.0.1:8080"}), rhConfig.Hosts, "hosts = [\"127.0.0.1:8080\"]")
	assert.Equal(t, 10, rhConfig.BatchSize, "batch_size = 10")
	assert.Equal(t, 5, rhConfig.RetryLimit, "retry_limit = 5")
	assert.Equal(t, time.Duration(2000000000), rhConfig.Timeout, "timeout = 2")
}

func TestValidateConfigError(t *testing.T) {
	config := `
hosts: ["127.0.0.1:8080"]
batch_size: -1
retry_limit: 5
timeout: 2
`
	c := common.MustNewConfigFrom(config)
	_, err := readConfig(c)
	assert.Errorf(t, err, "dbrainhub config params error")
}

func readConfig(cfg *common.Config) (*rainhubConfig, error) {
	c := rainhubConfig{}
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
