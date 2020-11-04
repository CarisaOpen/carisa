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

package object

import (
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/plugin"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

const locService = "object.service"

// Service implements CRUD operations for the plugin
type Service struct {
	cnt    *runtime.Container
	ext    *service.Extension
	crud   storage.CrudOperation
	plugin *plugin.Service
}

// NewService builds a plugin service
func NewService(cnt *runtime.Container, ext *service.Extension, crud storage.CrudOperation, plugin *plugin.Service) Service {
	return Service{
		cnt:    cnt,
		ext:    ext,
		crud:   crud,
		plugin: plugin,
	}
}

// Create creates a instance into of the repository.
// If the instance exists return false in the first param returned.
// If the plugin prototype doesn't exist return false in the second param returned.
// If the container doesn't exist return false in the third param returned.
func (s *Service) Create(inst *Instance) (bool, bool, bool, error) {
	inst.AutoID()

	found, err := s.plugin.Exists(inst.ProtoID)
	if err != nil {
		return false, false, false, err
	}
	if !found {
		return false, false, false, nil
	}

	created, foundc, err := s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, inst)
	return created, true, foundc, err
}

// Put creates or updates a instance into of the repository.
// If the instance exists return true in the first param returned otherwise return false.
// If the plugin prototype doesn't exist return false in the second param returned.
// The plugin is only checked when the instance exists.
// If the container doesn't exist return false in the third param returned.
func (s *Service) Put(inst *Instance) (bool, bool, bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	foundi, err := s.crud.Store().Exists(ctx, entity.ObjectKey(inst.ID))
	cancel()
	if err != nil {
		return false, false, false, err
	}
	if !foundi { // The plugin is only checked when the instance exists
		foundp, err := s.plugin.Exists(inst.ProtoID)
		if err != nil {
			return false, false, false, err
		}
		if !foundp {
			return false, false, false, nil
		}
	}

	updated, foundc, err := s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, inst)
	if err != nil {
		return false, true, false, err
	}
	return updated, true, foundc, nil
}

// Get gets the instance from storage
func (s *Service) Get(id xid.ID, inst *Instance) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, entity.ObjectKey(id), inst)
	cancel()
	return ok, err
}

// ListInstances lists queries depending ranges parameter.
// Look at service.List
func (s *Service) ListInstances(
	scheme string,
	id xid.ID,
	cat plugin.Category,
	name string,
	ranges bool,
	top int) ([]storage.Entity, error) {
	//
	return s.ext.List(
		entity.Key(scheme, id),
		strings.Concat(string(cat), name),
		ranges,
		top,
		func() storage.Entity { return &relation.PlatformInstance{} })
}
