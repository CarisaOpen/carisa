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

package factory

import (
	"testing"

	"go.etcd.io/etcd/integration"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/api/internal/runtime"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	cnf := runtime.Config{
		Server:     runtime.Server{Port: 8080},
		EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
	}

	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer cluster.Terminate(t)

	factory := build(cluster)

	assert.Equal(t, cnf, factory.Config, "Config")
	assert.NotNil(t, cnf, factory.Handlers.InstHandler, "Instance Handler")
}
