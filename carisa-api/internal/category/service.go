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

package category

import (
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/api/internal/service"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
)

const locService = "category.service"

// Service implements CRUD operations for the category
type Service struct {
	cnt  *runtime.Container
	ext  *service.Extension
	crud storage.CrudOperation
}

// NewService builds a category service
func NewService(cnt *runtime.Container, ext *service.Extension, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		ext:  ext,
		crud: crud,
	}
}

// Create creates a category into of the repository and links category and space or other category.
// If the category exists return false in the first param returned.
// If the space or category doesn't exist return false in the second param returned.
func (s *Service) Create(cat *Category) (bool, bool, error) {
	cat.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, cat)
}

// Put creates or updates a category into of the repository.
// If the category exists return true in the first param returned otherwise return false.
// If the space or cat doesn't exist return false in the second param returned.
func (s *Service) Put(cat *Category) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, cat)
}

// Get gets the category from storage
func (s *Service) Get(id xid.ID, cat *Category) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, id.String(), cat)
	cancel()
	return ok, err
}

// ListCategories lists categories depending ranges parameter.
// Look at service.List
func (s *Service) ListCategories(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(id, name, ranges, top, func() storage.Entity { return &relation.Hierarchy{} })
}
