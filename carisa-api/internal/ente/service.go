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

package ente

import (
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
)

const locService = "ente.service"

// Service implements CRUD operations for the ente domain
type Service struct {
	cnt  *runtime.Container
	crud storage.CrudOperation
}

// NewService builds a ente service
func NewService(cnt *runtime.Container, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		crud: crud,
	}
}

// Create creates a ente into of the repository and links ente and ente.
// If the ente exists return false in the first param returned.
// If the ente doesn't exist return false in the second param returned.
func (s *Service) Create(ente *Ente) (bool, bool, error) {
	ente.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, ente)
}

// Create creates or updates a ente into of the repository.
// If the ente exists return true in the first param returned otherwise return false.
// If the ente doesn't exist return false in the second param returned.
func (s *Service) Put(ente *Ente) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, ente)
}

// Get gets the ente from storage
func (s *Service) Get(id xid.ID, ente *Ente) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, id.String(), ente)
	cancel()
	return ok, err
}
