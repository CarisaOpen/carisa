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

package plugin

import (
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

const locService = "plugin.service"

// Service implements CRUD operations for the plugin
type Service struct {
	cnt  *runtime.Container
	ext  *service.Extension
	crud storage.CrudOperation
}

// NewService builds a plugin service
func NewService(cnt *runtime.Container, ext *service.Extension, crud storage.CrudOperation) Service {
	return Service{
		cnt:  cnt,
		ext:  ext,
		crud: crud,
	}
}

// Create creates a plugin into of the repository and links plugin and platform.
// If the plugin exists return false in the first param returned.
func (s *Service) Create(proto *Prototype) (bool, error) {
	proto.AutoID()
	created, _, err := s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, proto)
	return created, err
}

// Put creates or updates a plugin into of the repository.
// If the plugin exists return true in the first param returned otherwise return false.
func (s *Service) Put(proto *Prototype) (bool, error) {
	updated, _, err := s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, proto)
	return updated, err
}

// Get gets the plugin from storage
func (s *Service) Get(id xid.ID, proto *Prototype) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, entity.PluginKey(id), proto)
	cancel()
	return ok, err
}

// Exists checks if the plugin exists
func (s *Service) Exists(id xid.ID) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	found, err := s.crud.Store().Exists(ctx, entity.PluginKey(id))
	if err != nil {
		return false,
			s.cnt.Log.ErrWrap1(
				err,
				"checking if the plugin prototype exists",
				locService,
				logging.String("PrototypeID", id.String()))
	}
	cancel()
	return found, err
}

// ListPlugins lists the plugins depending 'ranges' parameter.
// Look at service.List
func (s *Service) ListPlugins(cat Category, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(
		storage.Virtual,
		strings.Concat(string(cat), name),
		ranges,
		top,
		func() storage.Entity { return &relation.PlatformPlugin{} })
}
