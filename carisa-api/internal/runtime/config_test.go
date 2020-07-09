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

	"github.com/carisa/pkg/logging"
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
		envC string
		cnf  Config
	}{
		{
			envC: "",
			cnf: Config{
				Server: Server{
					Port: 8080,
				},
				ZapConfig:  logging.ZapConfig{},
				EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
			},
		},
		{
			envC: `{
  "log": {
    "development": true, 
    "level": 2, 
    "encoding": "json"
  }, 
  "etcd": {
    "dialKeepAliveTime": 2, 
    "endpoints": [
      "server1", 
      "server2"
    ], 
    "dialTimeout": 1, 
    "dialKeepAliveTimeout": 3, 
    "requestTimeout": 4
  }, 
  "server": {
    "port": 1212
  }
}`,
			cnf: Config{
				Server: Server{
					Port: 1212,
				},
				ZapConfig: logging.ZapConfig{
					Development: true,
					Level:       2,
					Encoding:    "json",
				},
				EtcdConfig: storage.EtcdConfig{
					DialTimeout:          1,
					DialKeepAliveTime:    2,
					DialKeepAliveTimeout: 3,
					RequestTimeout:       4,
					Endpoints:            []string{"server1", "server2"},
				},
			},
		},
		{
			envC: `{
  "etcd": {
    "requestTimeout": 4, 
    "dialTimeout": 1, 
    "dialKeepAliveTimeout": 3
  }
}`,
			cnf: Config{
				Server: Server{
					Port: 8080,
				},
				ZapConfig: logging.ZapConfig{},
				EtcdConfig: storage.EtcdConfig{
					DialTimeout:          1,
					DialKeepAliveTime:    0,
					DialKeepAliveTimeout: 3,
					RequestTimeout:       4,
				},
			},
		},
	}
	for i, tt := range tests {
		if len(tt.envC) != 0 {
			_ = os.Setenv(envConfig, tt.envC)
		}
		cnf := LoadConfig()
		assert.Equalf(t, tt.cnf, cnf, "Configuration %v", i+1)
	}
}

func TestRuntime_LoadConfigPanic(t *testing.T) {
	err := os.Setenv(envConfig, "Panic")
	if assert.NoError(t, err) {
		assert.Panics(t, func() { LoadConfig() })
	}
}

func TestRuntime_StoreWithTimeout(t *testing.T) {
	c := Config{}
	ctx, _ := c.StoreWithTimeout()
	assert.NotNil(t, ctx, "Request timeout context")
}

func TestRuntime_ConfigToString(t *testing.T) {
	cnf := Config{
		Server: Server{
			Port: 8080,
		},
		ZapConfig:  logging.ZapConfig{},
		EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
	}
	assert.NotEmpty(t, cnf.String())
}
