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
	"time"

	"github.com/carisa/pkg/runtime"
)

const envConfig = "CARISA_CONTROLLER"

// Config defines the global information
type Config struct {
	// SplitterSelectorInSecs launches a process that prioritises the splitters
	// that will be chosen for the management of each Ente
	// according to the consumption of each one.
	SplitterSelectorInSecs time.Duration `json:"splitterSelectorInSecs,omitempty"`
	// SplitterMaxMemoryInMB is the maximum memory of work of a splitter
	SplitterMaxMemoryInMB uint32 `json:"splitterMaxMemory,omitempty"`

	runtime.CommonConfig
}

func (c *Config) Common() *runtime.CommonConfig {
	return &c.CommonConfig
}

// LoadConfig loads the configuration from environment variable
func LoadConfig() Config {
	cnf := Config{
		SplitterSelectorInSecs: 15,
		SplitterMaxMemoryInMB:  1024,
	}
	runtime.LoadConfig(envConfig, &cnf)
	return cnf
}
