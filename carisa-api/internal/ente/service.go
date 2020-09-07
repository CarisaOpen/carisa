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
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/api/internal/service"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
)

const locService = "ente.service"

// Service implements CRUD operations for the ente
type Service struct {
	cnt  *runtime.Container
	ext  *service.Extension
	crud storage.CrudOperation
}

// NewService builds a ente service
func NewService(cnt *runtime.Container, ext *service.Extension, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		ext:  ext,
		crud: crud,
	}
}

// Create creates a ente into of the repository and links ente and space.
// If the ente exists return false in the first param returned.
// If the space doesn't exist return false in the second param returned.
func (s *Service) Create(ente *Ente) (bool, bool, error) {
	ente.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, ente)
}

// Put creates or updates a ente into of the repository.
// If the ente exists return true in the first param returned otherwise return false.
// If the space doesn't exist return false in the second param returned.
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

// ListProps lists entes depending ranges parameter.
// Look at service.List
func (s *Service) ListProps(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(id, name, ranges, top, func() storage.Entity { return &relation.EnteEnteProp{} })
}

// CreateProp creates a property into of the repository and links ente property and ente.
// If the property exists return false in the first param returned.
// If the ente doesn't exist return false in the second param returned.
func (s *Service) CreateProp(prop *Prop) (bool, bool, error) {
	prop.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, prop)
}

// PutProp creates or updates a peroperty into of the repository.
// If the property exists return true in the first param returned otherwise return false.
// If the ente doesn't exist return false in the second param returned.
func (s *Service) PutProp(prop *Prop) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, prop)
}

// GetProp gets the property from storage
func (s *Service) GetProp(id xid.ID, prop *Prop) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, id.String(), prop)
	cancel()
	return ok, err
}
