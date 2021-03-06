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

	"github.com/carisa/internal/api/mock"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/runtime"
	"github.com/stretchr/testify/assert"

	pkgr "github.com/carisa/pkg/runtime"
)

func TestTemplate_Build(t *testing.T) {
	cnf := runtime.Config{
		Server: runtime.Server{Port: 8080},
		CommonConfig: pkgr.CommonConfig{
			EtcdConfig: storage.EtcdConfig{RequestTimeout: 10},
		},
	}

	sMock := mock.NewStorageFake(t)
	defer sMock.Close()

	factory := build(sMock)

	assert.Equal(t, cnf, factory.Config, "Config")
	assert.NotNil(t, cnf, factory.Echo, "Http")

	assert.NotNil(t, factory.Handlers.InstHandler, "Inst Handler")
	assert.NotNil(t, factory.Handlers.SpaceHandler, "Space Handler")
	assert.NotNil(t, factory.Handlers.EnteHandler, "Ente Handler")
	assert.NotNil(t, factory.Handlers.CategoryHandler, "Category Handler")
	assert.NotNil(t, factory.Handlers.PluginHandler, "Plugin Handler")
	assert.NotNil(t, factory.Handlers.ObjectHandler, "Object Handler")
}
