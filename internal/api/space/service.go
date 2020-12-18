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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

const locService = "space.service"

// Service implements CRUD operations for the space.Space
type Service struct {
	cnt  *runtime.Container
	ext  *service.Extension
	crud storage.CrudOperation
}

// NewService builds a space.Space service
func NewService(cnt *runtime.Container, ext *service.Extension, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		ext:  ext,
		crud: crud,
	}
}

// Create creates a space into of the repository and links instance.Instance and space.Space.
// If the space.Space exists return false in the first param returned.
// If the instance.Instance doesn't exist return false in the second param returned.
func (s *Service) Create(space *Space) (bool, bool, error) {
	space.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, space)
}

// Put creates or updates a space.Space into of the repository.
// If the space exists return true in the first param returned otherwise return false.
// If the instance.Instance doesn't exist return false in the second param returned.
func (s *Service) Put(space *Space) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, space)
}

// Get gets the space.Space from storage
func (s *Service) Get(id xid.ID, space *Space) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, entity.SpaceKey(id), space)
	cancel()
	return ok, err
}

// ListEntes lists entes depending 'ranges' parameter.
// Look at service.List
func (s *Service) ListEntes(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(
		entity.SpaceKey(id),
		strings.Concat(relation.SpaceEnteLn, name),
		ranges,
		top,
		func() storage.Entity { return &relation.SpaceEnte{} })
}

// ListCategories lists categories depending 'ranges' parameter.
// Look at service.List
func (s *Service) ListCategories(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(
		entity.SpaceKey(id),
		strings.Concat(relation.SpaceCatLn, name),
		ranges,
		top,
		func() storage.Entity { return &relation.SpaceCategory{} })
}
