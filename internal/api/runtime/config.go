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
	"strconv"

	"github.com/carisa/pkg/runtime"
	"github.com/carisa/pkg/strings"
)

const envConfig = "CARISA_API"

// Server describes the http configuration
type Server struct {
	Port int `json:"port"`
}

// Address returns address to connection server
func (s *Server) Address() string {
	return strings.Concat(":", strconv.Itoa(s.Port))
}

// Config defines the global information
type Config struct {
	Server `json:"server,omitempty"`
	runtime.CommonConfig
}

func (c *Config) Common() *runtime.CommonConfig {
	return &c.CommonConfig
}

// LoadConfig loads the configuration from environment variable
func LoadConfig() Config {
	cnf := Config{
		Server: Server{
			Port: 8080,
		},
	}
	runtime.LoadConfig(envConfig, &cnf)
	return cnf
}
