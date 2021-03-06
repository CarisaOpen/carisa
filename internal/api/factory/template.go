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
	"github.com/carisa/internal/api/http/handler"
	"github.com/carisa/internal/api/runtime"
	loge "github.com/carisa/pkg/http/echo"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
	"github.com/labstack/echo/v4"
)

const locBuild = "factory.build"

// Template builds the dependencies for the application
type Template struct {
	Config   runtime.Config
	Handlers handler.Handlers
	Echo     *echo.Echo

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

// Build builds the services, store, log, etc..
func Build() Template {
	return build(nil)
}

func build(mng storage.Integration /*for test*/) Template {
	cnf, cnt, store, e := servers(mng)
	srv := services(cnt, store)
	handlers := handlers(srv, cnt)
	cnt.Log.Info1("http server started", locBuild, logging.String("address", cnf.Server.Address()))

	return Template{
		Config:   cnf,
		Handlers: handlers,
		Echo:     e,
		store:    store,
		cnt:      cnt,
	}
}

func servers(mng storage.Integration) (runtime.Config, *runtime.Container, storage.CRUD, *echo.Echo) {
	cnf := runtime.LoadConfig()
	log, zLog := logging.NewZapLogger(cnf.ZapConfig)
	log.Info1("loaded configuration", locBuild, logging.String("config", cnf.String()))

	log.Info1("initializing http server", locBuild, logging.String("address", cnf.Server.Address()))
	e := echo.New()
	e.Logger = loge.NewLogging("echo", loge.ConvertLevel(log.Level()), zLog)

	cnt := runtime.NewContainer(cnf, log)

	log.Info1("starting etcd client", locBuild, logging.String("endpoints", cnf.EPSString()))
	var store storage.CRUD
	if mng != nil {
		store = mng.Store()
	} else {
		store = storage.NewEtcdConfig(cnf.EtcdConfig)
	}
	return cnf, cnt, store, e
}

func services(cnt *runtime.Container, store storage.CRUD) service {
	cnt.Log.Info("configuring services", locBuild)
	srv := configService(cnt, store)
	return srv
}

func handlers(srv service, cnt *runtime.Container) handler.Handlers {
	cnt.Log.Info("configuring http handlers", locBuild)
	return handler.Handlers{
		InstHandler:     handler.NewInstanceHandle(srv.instanceSrv, cnt),
		SpaceHandler:    handler.NewSpaceHandle(srv.spaceSrv, cnt),
		EnteHandler:     handler.NewEnteHandle(srv.enteSrv, cnt),
		CategoryHandler: handler.NewCatHandle(srv.catSrv, cnt),
		PluginHandler:   handler.NewPluginHandle(srv.pluginSrv, cnt),
		ObjectHandler:   handler.NewObjectHandle(srv.objectSrv, cnt),
	}
}
