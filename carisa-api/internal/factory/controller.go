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
	"github.com/carisa/api/internal/http/handler"
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
)

const locBuild = "factory.build"

// Controller builds the application flow
type Controller struct {
	Config   runtime.Config
	Handlers handler.Handlers
	store    storage.CRUD
	cnt      runtime.Container
}

// Build builds the services, store, log, etc..
func Build() Controller {
	return build(nil)
}

func build(mng storage.Integration /*for test*/) Controller {
	cnf, cnt, store := servers(mng)
	srv := services(cnt, store)
	instHandler := handlers(srv, cnt)

	return Controller{
		Config: cnf,
		Handlers: handler.Handlers{
			InstHandler: instHandler,
		},
		store: store,
		cnt:   cnt,
	}
}

func servers(mng storage.Integration) (runtime.Config, runtime.Container, storage.CRUD) {
	cnf := runtime.LoadConfig()
	log := logging.NewZapLogger(cnf.ZapConfig)
	log.Info("loaded configuration", locBuild, logging.String("config", cnf.String()))
	cnt := runtime.NewContainer(cnf, log)
	log.Info("starting etcd client", locBuild, logging.String("endpoints", cnf.EPSString()))
	var store storage.CRUD
	if mng != nil {
		store = mng.Store()
	} else {
		store = storage.NewEtcdConfig(cnf.EtcdConfig)
	}
	return cnf, cnt, store
}

func services(cnt runtime.Container, store storage.CRUD) service {
	cnt.Log.Info("configuring services", locBuild)
	srv := configService(cnt, store)
	return srv
}

func handlers(srv service, cnt runtime.Container) handler.Instance {
	cnt.Log.Info("configuring http handlers", locBuild)
	instHandler := handler.NewInstanceHandl(srv.instanceSrv, cnt)
	return instHandler
}

// Close closes all connections
func (c *Controller) Close() {
	c.cnt.Log.Info("closing connections. bye", locBuild)
	if err := c.store.Close(); err != nil {
		_ = c.cnt.Log.ErrWrap(err, "closing storage", "Factory.Close")
	}
}
