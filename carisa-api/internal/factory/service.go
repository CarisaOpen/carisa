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
	"github.com/carisa/api/internal/category"
	"github.com/carisa/api/internal/ente"
	"github.com/carisa/api/internal/instance"
	"github.com/carisa/api/internal/runtime"
	srv "github.com/carisa/api/internal/service"
	"github.com/carisa/api/internal/space"
	"github.com/carisa/pkg/storage"
)

// Service configures all transversal services for API
type service struct {
	instanceSrv instance.Service
	spaceSrv    space.Service
	enteSrv     ente.Service
	catSrv      category.Service
}

// configService builds the services
func configService(cnt *runtime.Container, store storage.CRUD) service {
	crud := storage.NewCrudOperation(store, cnt.Log, storage.NewTxn)
	ext := srv.NewExt(cnt, store)
	return service{
		instanceSrv: instance.NewService(cnt, ext, crud),
		spaceSrv:    space.NewService(cnt, ext, crud),
		enteSrv:     ente.NewService(cnt, ext, crud),
		catSrv:      category.NewService(cnt, ext, crud),
	}
}
