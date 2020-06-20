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
	"github.com/carisa/api/internal/instance"
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/storage"
)

// Service configures all transversal services for API
type service struct {
	instanceSrv instance.Service
}

// configService builds the services
func configService(cnt runtime.Container, store storage.CRUD) service {
	return service{
		instanceSrv: instance.NewService(cnt, store),
	}
}
