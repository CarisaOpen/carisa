/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing,
 *   software  distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and  limitations under the License.
 */

package service

import (
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
)

type Extension struct {
	cnt  *runtime.Container
	crud storage.CRUD
}

// NewExt builds a extension for services
func NewExt(cnt *runtime.Container, crud storage.CRUD) *Extension {
	return &Extension{
		cnt:  cnt,
		crud: crud,
	}
}

// List lists entities depending ranges parameter.
// If ranges is equal to true is filtered by entity which name is greater than name parameter
// If ranges is equal to false is filtered by entity which name starts by name parameter
func (e *Extension) List(id string, name string, ranges bool, top int, empty func() storage.Entity) ([]storage.Entity, error) {
	var list []storage.Entity
	var err error

	ctx, cancel := e.cnt.StoreWithTimeout()
	if ranges {
		list, err = e.crud.Range(ctx, strings.Concat(id, name), id, top, empty)
	} else {
		list, err = e.crud.StartKey(ctx, strings.Concat(id, name), top, empty)
	}
	cancel()

	return list, err
}
