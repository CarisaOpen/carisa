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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// The thinks of space.Space.
// The ente.Ente are the items of space.Space to trace, count, measure, etc.
type Ente struct {
	entity.Descriptor
	SpaceID xid.ID `json:"spaceId"` // space.Space container
	CatID   xid.ID `json:"-"`       // Is used temporarily to connect the entity and the category.Category.
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
	return entity.EnteKey(e.ID)
}

func (e *Ente) Nominative() entity.Descriptor {
	return e.Descriptor
}

func (e *Ente) RelName() string {
	return e.Name
}

// ParentKey gets the Space ID
func (e *Ente) ParentKey() string {
	parentID := entity.CategoryKey(e.CatID)
	if e.CatID.IsNil() {
		parentID = entity.SpaceKey(e.SpaceID)
	}
	return parentID
}

// Link gets the link between instance.Instance and prop
func (e *Ente) Link() storage.Entity {
	cat := true
	if e.CatID.IsNil() {
		cat = false
	}
	return e.link(cat, e.ParentKey())
}

func (e *Ente) LinkName() string {
	if e.CatID.IsNil() {
		return relation.SpaceEnteLn
	}
	return relation.CatEnteLn
}

func (e *Ente) ReLink(dlr storage.DLRel) storage.Entity {
	cat := true
	if dlr.Type == relation.SpaceEnteLn {
		cat = false
	}
	return e.link(cat, dlr.ParentID)
}

func (e *Ente) link(cat bool, parentID string) storage.Entity {
	if cat {
		return &relation.Hierarchy{
			ID:       strings.Concat(parentID, e.Name, e.Key()),
			Name:     e.Name,
			LinkID:   e.ID.String(),
			Category: false,
		}
	}
	return &relation.SpaceEnte{
		ID:     strings.Concat(parentID, relation.SpaceEnteLn, e.Name, e.Key()),
		Name:   e.Name,
		EnteID: e.ID.String(),
	}
}

func (e *Ente) Empty() storage.EntityRelation {
	return &Ente{}
}

// The prop are properties contains the fields
type Prop struct {
	entity.Descriptor
	EnteID    xid.ID          `json:"enteId"` // Ente container
	Type      entity.TypeProp `json:"type"`
	CatPropID xid.ID          `json:"-"` // Is used temporarily to connect the property and the category property.
}

func NewProp() Prop {
	return Prop{
		Descriptor: entity.NewDescriptor(),
		Type:       entity.Integer,
	}
}

func (e *Prop) GetType() entity.TypeProp {
	return e.Type
}

func (e *Prop) ToString() string {
	return strings.Concat("ente-property: ID:", e.Key(), ", name:", e.Name)
}

func (e *Prop) Key() string {
	return entity.EntePropKey(e.ID)
}

func (e *Prop) Nominative() entity.Descriptor {
	return e.Descriptor
}

func (e *Prop) RelName() string {
	return e.Name
}

// ParentKey gets the Ente or category.Category ID
func (e *Prop) ParentKey() string {
	parentID := entity.CatPropKey(e.CatPropID)
	if e.CatPropID.IsNil() {
		parentID = entity.EnteKey(e.EnteID)
	}
	return parentID
}

func (e *Prop) Link() storage.Entity {
	cat := true
	if e.CatPropID.IsNil() {
		cat = false
	}
	return e.link(cat, e.ParentKey())
}

func (e *Prop) LinkName() string {
	if e.CatPropID.IsNil() {
		return relation.EntePropLn
	}
	return relation.CatPropPropLn
}

func (e *Prop) ReLink(dlr storage.DLRel) storage.Entity {
	cat := true
	if dlr.Type == relation.EntePropLn {
		cat = false
	}
	return e.link(cat, dlr.ParentID)
}

func (e *Prop) link(cat bool, parentID string) storage.Entity {
	if cat {
		return &relation.CatPropProp{
			ID:       strings.Concat(parentID, e.Name, e.Key()),
			Name:     e.Name,
			PropID:   e.ID.String(),
			Category: false,
		}
	}
	return &relation.EnteProp{
		ID:         strings.Concat(parentID, relation.EntePropLn, e.Name, e.Key()),
		Name:       e.Name,
		EntePropID: e.ID.String(),
	}
}

func (e *Prop) Empty() storage.EntityRelation {
	return &Prop{}
}
