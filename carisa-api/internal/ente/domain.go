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
	"github.com/carisa/api/internal/entity"
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// The thinks of spaces.
// The ente are the items of spaces to trace, count, measure, etc.
type Ente struct {
	entity.Descriptor
	SpaceID xid.ID `json:"spaceId"` // Space container
}

func New() Ente {
	return Ente{
		Descriptor: entity.NewDescriptor(),
	}
}

func (e *Ente) ToString() string {
	return strings.Concat("ente: ID:", e.Key(), ", name:", e.Name)
}

func (e *Ente) Key() string {
	return e.ID.String()
}

func (e *Ente) RelKey() string {
	return strings.Concat(e.SpaceID.String(), e.Name, e.Key())
}

func (e *Ente) RelName() string {
	return e.Name
}

// ParentKey gets the Space ID
func (e *Ente) ParentKey() string {
	return e.SpaceID.String()
}

// Link gets the link between instance and ente
func (e *Ente) Link() storage.Entity {
	return &relation.SpaceEnte{
		ID:     e.RelKey(),
		Name:   e.Name,
		EnteID: e.ID.String(),
	}
}

func (e *Ente) Empty() storage.EntityRelation {
	return &Ente{}
}
