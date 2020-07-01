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
	"testing"

	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestRuntime_NewContainer(t *testing.T) {
	cnf := Config{
		Server:     Server{Port: 8080},
		EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
	}

	log := logging.NewZapLogger(cnf.ZapConfig)
	sMock := storage.NewEctdIntegra(t)
	defer sMock.Close()

	ctn := NewContainer(cnf, log)

	assert.Equal(t, cnf, ctn.Config, "Config")
	assert.NotNil(t, ctn.Log, "Log")
	assert.NotNil(t, ctn.NewTxn(sMock.Store()), "Transaction")
}
