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

package instance

import (
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/storage"
)

const locService = "instance.service"

// Service implements CRUD operations for the instance domain service
type Service struct {
	cnt  *runtime.Container
	crud storage.CrudOperation
}

// NewService builds a instance service
func NewService(cnt *runtime.Container, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		crud: crud,
	}
}

// Create creates a instance into of the repository
// If the instance exists returns false
func (s *Service) Create(inst *Instance) (bool, error) {
	inst.AutoID()
	return s.crud.Create(locService, s.cnt.StoreWithTimeout, inst)
}

// Put creates or updates depending of if exists the instance into storage
// If the instance is updated return true
func (s *Service) Put(inst *Instance) (bool, error) {
	return s.crud.Put(locService, s.cnt.StoreWithTimeout, inst)
}
