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
	"strconv"
	"time"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/pkg/logging"

	"github.com/carisa/pkg/strings"
)

const envConfig = "CARISA_API"

// Server describes the http configuration
type Server struct {
	Port int `json:"port"`
}

// Address return address to connection server
func (s *Server) Address() string {
	return strings.Concat(":", strconv.Itoa(s.Port))
}

// Config defines the global information
type Config struct {
	Server             `json:"server,omitempty"`
	logging.ZapConfig  `json:"log,omitempty"`
	storage.EtcdConfig `json:"etcd,omitempty"`
}

// LoadConfig loads the configuration from environment variable
func LoadConfig() Config {
	env := os.Getenv(envConfig)

	cnf := Config{
		Server: Server{
			Port: 8080,
		},
		EtcdConfig: storage.EtcdConfig{},
		ZapConfig:  logging.ZapConfig{},
	}

	if len(env) != 0 {
		if err := json.Unmarshal([]byte(env), &cnf); err != nil {
			panic(strings.Concat("configuration environment variable cannot be loaded: ", err.Error()))
		}
	}
	// Default values not treated in the factories
	if cnf.RequestTimeout == 0 {
		cnf.RequestTimeout = 10
	}

	return cnf
}

// StoreWithTimeout creates the timeout context with the value Store.RequestTimeout
func (c *Config) StoreWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(c.RequestTimeout)*time.Second)
}

func (c *Config) String() string {
	r, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(r)
}
