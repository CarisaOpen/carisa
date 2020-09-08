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
	"github.com/carisa/api/internal/ente"
	"github.com/carisa/api/internal/entity"
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// Category represents a hierarchy to organize the entes and how the information spreads between the different categories
// Each category references to cat properties or properties of other categories.
// If it is Root, the parent is a space.
type Category struct {
	entity.Descriptor
	ParentID xid.ID `json:"parentId,omitempty"` // Space or category container
	Root     bool   `json:"root"`
}

func New() Category {
	return Category{
		Descriptor: entity.NewDescriptor(),
	}
}

func (c *Category) ToString() string {
	return strings.Concat("category: ID:", c.Key(), ", name:", c.Name)
}

func (c *Category) Key() string {
	return c.ID.String()
}

func (c *Category) Nominative() entity.Descriptor {
	return c.Descriptor
}

func (c *Category) RelKey() string {
	if c.Root {
		return strings.Concat(c.ParentID.String(), "S", c.Name, c.Key())
	}
	return strings.Concat(c.ParentID.String(), "C", c.Name, c.Key())
}

func (c *Category) RelName() string {
	return c.Name
}

func (c *Category) ParentKey() string {
	return c.ParentID.String()
}

func (c *Category) SetParentKey(value string) error {
	id, err := xid.FromString(value)
	if err != nil {
		return err
	}
	c.ParentID = id
	return nil
}

// Link gets the link between category and space, if the Root field is true
// If the Root field is false gets the link between category and others category or cat
func (c *Category) Link() storage.Entity {
	if c.Root {
		return &relation.SpaceCategory{
			ID:    c.RelKey(),
			Name:  c.Name,
			CatID: c.ID.String(),
		}
	}
	return &relation.Hierarchy{
		ID:       c.RelKey(),
		Name:     c.Name,
		LinkID:   c.ID.String(),
		Category: true,
	}
}

func (c *Category) Empty() storage.EntityRelation {
	return &Category{}
}

// Prop are the properties of category. Each property represents a generalization
// of the properties of the entes or of the properties of the categories
// All of them have the same type. The type is assigned when is linked the first property so ente as the category
type Prop struct {
	entity.Descriptor
	CatID xid.ID        `json:"categoryId"` // Category container
	Type  ente.TypeProp `json:"type"`
}

func NewProp() Prop {
	return Prop{
		Descriptor: entity.NewDescriptor(),
		Type:       ente.None,
	}
}

func (c *Prop) ToString() string {
	return strings.Concat("category-property: ID:", c.Key(), ", name:", c.Name)
}

func (c *Prop) Key() string {
	return c.ID.String()
}

func (c *Prop) Nominative() entity.Descriptor {
	return c.Descriptor
}

func (c *Prop) RelKey() string {
	return strings.Concat(c.CatID.String(), "P", c.Name, c.Key())
}

func (c *Prop) RelName() string {
	return c.Name
}

// ParentKey gets the Ente ID
func (c *Prop) ParentKey() string {
	return c.CatID.String()
}

func (c *Prop) SetParentKey(value string) error {
	id, err := xid.FromString(value)
	if err != nil {
		return err
	}
	c.CatID = id
	return nil
}

// Link gets the link between Category and her properties
func (c *Prop) Link() storage.Entity {
	return &relation.CategoryProp{
		ID:        c.RelKey(),
		Name:      c.Name,
		CatPropID: c.ID.String(),
	}
}

func (c *Prop) Empty() storage.EntityRelation {
	return &Prop{}
}
