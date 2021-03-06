/*
 *  Copyright 2019-2022 the original author or authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package factory

import (
	"github.com/carisa/internal/splitter/runtime"
	"github.com/carisa/internal/splitter/service"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
)

const locBuild = "factory.build"

// Template builds the dependencies for the application
type Template struct {
	Controller service.Controller

	store storage.CRUD
	cnt   *runtime.Container
}

func (c *Template) Close() {
	const loc = "factory.close"
	c.cnt.Log.Info("closing connections", loc)
	if err := c.store.Close(); err != nil {
		c.cnt.Log.ErrorE(err, loc)
	} else {
		c.cnt.Log.Info("closed connections", loc)
	}
}

// Build builds the controllers, store, log, etc..
func Build() Template {
	return build(nil)
}

func build(mng storage.Integration /*for test*/) Template {
	cnt, store := servers(mng)

	return Template{
		Controller: service.NewController(cnt, store),
		store:      store,
		cnt:        cnt,
	}
}

func servers(mng storage.Integration) (*runtime.Container, storage.CRUD) {
	cnf := runtime.LoadConfig()
	log, _ := logging.NewZapLogger(cnf.ZapConfig)
	log.Info1("loaded configuration", locBuild, logging.String("config", cnf.String()))

	cnt := runtime.NewContainer(cnf, storage.NewTxn, log)

	log.Info1("starting etcd client", locBuild, logging.String("endpoints", cnf.EPSString()))
	var store storage.CRUD
	if mng != nil {
		store = mng.Store()
	} else {
		store = storage.NewEtcdConfig(cnf.EtcdConfig)
	}
	return cnt, store
}
