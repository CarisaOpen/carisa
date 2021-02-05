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
	"github.com/carisa/pkg/runtime"
	"github.com/rs/xid"
)

const envConfig = "CARISA_SPLITTER"

// Server describes the information about splitter server
type Server struct {
	Name xid.ID // Unique identifier for each splitter
}

func newServer() Server {
	return Server{
		Name: xid.New(),
	}
}

// Config defines the global information
type Config struct {
	Server `json:"-"`
	runtime.CommonConfig
}

func (c *Config) Common() *runtime.CommonConfig {
	return &c.CommonConfig
}

// LoadConfig loads the configuration from environment variable
func LoadConfig() Config {
	cnf := Config{
		Server: newServer(),
	}
	runtime.LoadConfig(envConfig, &cnf)
	return cnf
}
