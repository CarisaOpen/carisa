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

package space

import (
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/storage"
)

const locService = "space.service"

// Service implements CRUD operations for the space domain
type Service struct {
	cnt  *runtime.Container
	crud storage.CrudOperation
}

// NewService builds a space service
func NewService(cnt *runtime.Container, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		crud: crud,
	}
}

// Create creates a space into of the repository and links instance and space.
// If the space exists return false in the first param returned.
// If the instance doesn't exist return false in the second param returned.
func (s *Service) Create(space *Space) (bool, bool, error) {
	space.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, space)
}

// Create creates or updates a space into of the repository.
// If the space exists return true in the first param returned otherwise return false.
// If the instance doesn't exist return false in the second param returned.
func (s *Service) Put(space *Space) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, space)
}
