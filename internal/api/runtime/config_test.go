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
 * software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package runtime

import (
	"os"
	"testing"

	"github.com/carisa/pkg/runtime"
	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestServer_Address(t *testing.T) {
	s := Server{
		Port: 1212,
	}
	assert.Equal(t, ":1212", s.Address())
}

func TestRuntime_LoadConfig(t *testing.T) {
	tests := []struct {
		name string
		envC string
		cnf  Config
	}{
		{
			name: "Default configuration",
			envC: "",
			cnf: Config{
				Server: Server{
					Port: 8080,
				},
				CommonConfig: runtime.CommonConfig{
					EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
				},
			},
		},
		{
			name: "Server configuration",
			envC: `{
  "server": {
    "port": 1212
  }
}`,
			cnf: Config{
				Server: Server{
					Port: 1212,
				},
				CommonConfig: runtime.CommonConfig{
					EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
				},
			},
		},
	}
	for _, tt := range tests {
		if len(tt.envC) != 0 {
			_ = os.Setenv(envConfig, tt.envC)
		}
		cnf := LoadConfig()
		assert.Equalf(t, tt.cnf, cnf, tt.name)
	}
}
