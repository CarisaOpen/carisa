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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// Category represents a hierarchy to organize the ente.Ente and how the information spreads between the different categories
// Each Category references to cat properties or properties of other categories.
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
	return entity.CategoryKey(c.ID)
}

func (c *Category) Nominative() entity.Descriptor {
	return c.Descriptor
}

func (c *Category) RelName() string {
	return c.Name
}

func (c *Category) ParentKey() string {
	if c.Root {
		return entity.SpaceKey(c.ParentID)
	}
	return entity.CategoryKey(c.ParentID)
}

// Link gets the link between Category and space.Space if the Root field is true
// If the Root field is false gets the link between category and others Category
func (c *Category) Link() storage.Entity {
	return c.link(c.Root, c.ParentKey())
}

func (c *Category) LinkName() string {
	if c.Root {
		return relation.SpaceCatLn
	}
	return relation.CatCatLn
}

func (c *Category) ReLink(dlr storage.DLRel) storage.Entity {
	root := true
	if dlr.Type == relation.CatCatLn {
		root = false
	}
	return c.link(root, dlr.ParentID)
}

func (c *Category) link(root bool, parentID string) storage.Entity {
	if root {
		return &relation.SpaceCategory{
			ID:    strings.Concat(parentID, relation.SpaceCatLn, c.Name, c.Key()),
			Name:  c.Name,
			CatID: c.ID.String(),
		}
	}
	return &relation.Hierarchy{
		ID:       strings.Concat(parentID, c.Name, c.Key()),
		Name:     c.Name,
		LinkID:   c.ID.String(),
		Category: true,
	}
}

func (c *Category) Empty() storage.EntityRelation {
	return &Category{}
}

// Prop are the properties of Category. Each property represents a generalization
// of the properties of the ente.Ente or of the properties of the categories
// Each Category property can inherit properties of child (Category or ente.Ente)
// All of them have the same type. The type is assigned when is linked the first property so ente.Ente as the Category
type Prop struct {
	entity.Descriptor
	CatID     xid.ID          `json:"categoryId"` // Category container
	Type      entity.TypeProp `json:"type"`
	catPropID xid.ID          // Is used temporarily to connect the property and the Category property.
}

func NewProp() Prop {
	return Prop{
		Descriptor: entity.NewDescriptor(),
		Type:       entity.None,
	}
}

func (c *Prop) GetType() entity.TypeProp {
	return c.Type
}

func (c *Prop) ToString() string {
	return strings.Concat("category-property: ID:", c.Key(), ", name:", c.Name)
}

func (c *Prop) Key() string {
	return entity.CatPropKey(c.ID)
}

func (c *Prop) Nominative() entity.Descriptor {
	return c.Descriptor
}

func (c *Prop) RelName() string {
	return c.Name
}

func (c *Prop) ParentKey() string {
	parentID := entity.CatPropKey(c.catPropID)
	if c.catPropID.IsNil() {
		parentID = entity.CategoryKey(c.CatID)
	}
	return parentID
}

func (c *Prop) Link() storage.Entity {
	catProp := true
	if c.catPropID.IsNil() {
		catProp = false
	}
	return c.link(catProp, c.ParentKey())
}

func (c *Prop) LinkName() string {
	if c.catPropID.IsNil() {
		return relation.CatPropLn
	}
	return relation.CatPropPropLn
}

func (c *Prop) ReLink(dlr storage.DLRel) storage.Entity {
	catProp := true
	if dlr.Type == relation.CatPropLn {
		catProp = false
	}
	return c.link(catProp, dlr.ParentID)
}

func (c *Prop) link(catProp bool, parentID string) storage.Entity {
	if catProp {
		return &relation.CatPropProp{
			ID:       strings.Concat(parentID, c.Name, c.Key()),
			Name:     c.Name,
			PropID:   c.ID.String(),
			Category: true,
		}
	}
	return &relation.CategoryProp{
		ID:        strings.Concat(parentID, relation.CatPropLn, c.Name, c.Key()),
		Name:      c.Name,
		CatPropID: c.ID.String(),
	}
}

func (c *Prop) Empty() storage.EntityRelation {
	return &Prop{}
}
