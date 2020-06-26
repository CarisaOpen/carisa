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

// Factory builds the application flow
type Factory struct {
	Config   runtime.Config
	Handlers handler.Handlers
}

// Build builds the services, store, log, etc..
func Build() Factory {
	return build(nil)
}

func build(mng storage.Integration /*for test*/) Factory {
	cnf, cnt, store := servers(mng)
	srv := services(cnt, store)
	instHandler := handlers(srv, cnt)

	return Factory{
		Config: cnf,
		Handlers: handler.Handlers{
			InstHandler: instHandler,
		},
	}
}

func servers(mng storage.Integration) (runtime.Config, runtime.Container, storage.CRUD) {
	cnf := runtime.LoadConfig()
	log := logging.NewZapLogger(cnf.ZapConfig)
	cnt := runtime.Container{
		Config: cnf,
		Log:    log,
	}
	var store storage.CRUD
	if mng != nil {
		store = mng.Store()
	} else {
		store = storage.NewEtcdConfig(cnf.EtcdConfig)
	}
	return cnf, cnt, store
}

func services(cnt runtime.Container, store storage.CRUD) service {
	srv := configService(cnt, store)
	return srv
}

func handlers(srv service, cnt runtime.Container) handler.Instance {
	instHandler := handler.NewInstanceHandl(srv.instanceSrv, cnt)
	return instHandler
}
