/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package runtime

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/pkg/logging"

	"github.com/carisa/pkg/strings"
)

// Config is the interface to get common config
type Config interface {
	Common() *CommonConfig
}

// Config defines the configuration by default
type CommonConfig struct {
	logging.ZapConfig  `json:"log,omitempty"`
	storage.EtcdConfig `json:"etcd,omitempty"`
}

// LoadConfig loads the configuration from environment variable
func LoadConfig(envConfig string, cnf Config) {
	env := os.Getenv(envConfig)

	if len(env) != 0 {
		if err := json.Unmarshal([]byte(env), &cnf); err != nil {
			panic(strings.Concat("configuration environment variable cannot be loaded: ", err.Error()))
		}
	}
	dCnf := cnf.Common()
	// Common values not treated in the factories
	if dCnf.RequestTimeout == 0 {
		dCnf.RequestTimeout = 10
	}
}

// StoreWithTimeout creates the timeout context with the value Store.RequestTimeout
func (c *CommonConfig) StoreWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(c.RequestTimeout)*time.Second)
}

func (c *CommonConfig) String() string {
	r, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(r)
}
