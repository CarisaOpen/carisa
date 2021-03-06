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
	"testing"

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"

	"github.com/rs/xid"

	"github.com/carisa/pkg/strings"

	"github.com/stretchr/testify/assert"
)

func TestCategory_ToString(t *testing.T) {
	c := New()
	assert.Equal(t, strings.Concat("category: ID:", c.Key(), ", name:", c.Name), c.ToString())
}

func TestCategory_Key(t *testing.T) {
	c := New()
	assert.Equal(t, entity.CategoryKey(c.ID), c.Key())
}

func TestCategory_Nominative(t *testing.T) {
	c := Category{}
	assert.Equal(t, entity.Descriptor{}, c.Nominative())
}

func TestCategory_ParentKey(t *testing.T) {
	c := New()
	c.ParentID = xid.New()
	assert.Equal(t, entity.CategoryKey(c.ParentID), c.ParentKey())
}

func TestCategory_Empty(t *testing.T) {
	c := New()
	assert.Equal(t, &Category{}, c.Empty())
}

func TestCategory_Link(t *testing.T) {
	c := New()
	c.ParentID = xid.New()

	tests := []struct {
		name string
		root bool
		link storage.Entity
	}{
		{
			name: "Category with space",
			root: true,
			link: &relation.SpaceCategory{
				ID:    strings.Concat(entity.SpaceKey(c.ParentID), relation.SpaceCatLn, c.Name, c.Key()),
				Name:  c.Name,
				CatID: c.ID.String(),
			},
		},
		{
			name: "Category with others Category",
			root: false,
			link: &relation.Hierarchy{
				ID:       strings.Concat(entity.CategoryKey(c.ParentID), c.Name, c.Key()),
				Name:     c.Name,
				Category: true,
				LinkID:   c.ID.String(),
			},
		},
	}

	for _, tt := range tests {
		c.Root = tt.root
		assert.Equal(t, tt.link, c.Link(), tt.name)
	}
}

func TestCategory_LinkName(t *testing.T) {
	c := New()

	tests := []struct {
		name  string
		root  bool
		typen string
	}{
		{
			name:  "Space -> Category",
			root:  true,
			typen: relation.SpaceCatLn,
		},
		{
			name:  "Category -> Category",
			root:  false,
			typen: relation.CatCatLn,
		},
	}

	for _, tt := range tests {
		c.Root = tt.root
		assert.Equal(t, tt.typen, c.LinkName(), tt.name)
	}
}

func TestCategory_ReLink(t *testing.T) {
	c := New()
	parentID := xid.New().String()

	tests := []struct {
		name string
		tn   string
		link storage.Entity
	}{
		{
			name: "Category with space",
			tn:   relation.SpaceCatLn,
			link: &relation.SpaceCategory{
				ID:    strings.Concat(parentID, relation.SpaceCatLn, c.Name, c.Key()),
				Name:  c.Name,
				CatID: c.ID.String(),
			},
		},
		{
			name: "Category with others Category",
			tn:   relation.CatCatLn,
			link: &relation.Hierarchy{
				ID:       strings.Concat(parentID, c.Name, c.Key()),
				Name:     c.Name,
				Category: true,
				LinkID:   c.ID.String(),
			},
		},
	}

	for _, tt := range tests {
		dlr := storage.DLRel{
			ParentID: parentID,
			Type:     tt.tn,
		}
		assert.Equal(t, tt.link, c.ReLink(dlr), tt.name)
	}
}

func TestCategoryCatProp_Field(t *testing.T) {
	c := NewProp()
	assert.Equal(t, c.GetType(), entity.None)
}

func TestCategoryCatProp_ToString(t *testing.T) {
	c := NewProp()
	assert.Equal(t, strings.Concat("category-property: ID:", c.Key(), ", name:", c.Name), c.ToString())
}

func TestCategoryCatProp_Key(t *testing.T) {
	c := NewProp()
	assert.Equal(t, entity.CatPropKey(c.ID), c.Key())
}

func TestCategoryCatProp_Nominative(t *testing.T) {
	c := Prop{}
	assert.Equal(t, entity.Descriptor{}, c.Nominative())
}

func TestCategoryCatProp_ParentKey(t *testing.T) {
	c := NewProp()
	c.CatID = xid.New()
	assert.Equal(t, entity.CategoryKey(c.CatID), c.ParentKey())
}

func TestCategoryCatProp_Empty(t *testing.T) {
	c := NewProp()
	assert.Equal(t, &Prop{}, c.Empty())
}

func TestCategoryCatProp_Link(t *testing.T) {
	propID := xid.New()
	parentID := xid.New()
	name := "namel"

	tests := []struct {
		name string
		prop Prop
		link storage.Entity
	}{
		{
			name: "Category -> Property",
			prop: Prop{
				Descriptor: entity.Descriptor{
					ID:   propID,
					Name: "namel",
				},
				CatID: parentID,
			},
			link: &relation.CategoryProp{
				ID:        strings.Concat(entity.CategoryKey(parentID), relation.CatPropLn, name, entity.CatPropKey(propID)),
				Name:      name,
				CatPropID: propID.String(),
			},
		},
		{
			name: "Category property -> Category property",
			prop: Prop{
				Descriptor: entity.Descriptor{
					ID:   propID,
					Name: "namel",
				},
				catPropID: parentID,
			},
			link: &relation.CatPropProp{
				ID:       strings.Concat(entity.CatPropKey(parentID), name, entity.CatPropKey(propID)),
				Name:     name,
				PropID:   propID.String(),
				Category: true,
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.link, tt.prop.Link(), tt.name)
	}
}

func TestCategoryCatProp_LinkName(t *testing.T) {
	p := NewProp()

	tests := []struct {
		name      string
		catPropID xid.ID
		typen     string
	}{
		{
			name:  "Category -> Property",
			typen: relation.CatPropLn,
		},
		{
			name:      "Category property -> Category Property",
			catPropID: xid.New(),
			typen:     relation.CatPropPropLn,
		},
	}

	for _, tt := range tests {
		p.catPropID = tt.catPropID
		assert.Equal(t, tt.typen, p.LinkName(), tt.name)
	}
}

func TestCategoryCatProp_ReLink(t *testing.T) {
	p := NewProp()
	parentID := xid.New().String()

	tests := []struct {
		name string
		tn   string
		link storage.Entity
	}{
		{
			name: "Category -> Property",
			tn:   relation.CatPropLn,
			link: &relation.CategoryProp{
				ID:        strings.Concat(parentID, relation.CatPropLn, p.Name, p.Key()),
				Name:      p.Name,
				CatPropID: p.ID.String(),
			},
		},
		{
			name: "Category property -> Category Property",
			tn:   relation.CatPropPropLn,
			link: &relation.CatPropProp{
				ID:       strings.Concat(parentID, p.Name, p.Key()),
				Name:     p.Name,
				Category: true,
				PropID:   p.ID.String(),
			},
		},
	}

	for _, tt := range tests {
		dlr := storage.DLRel{
			ParentID: parentID,
			Type:     tt.tn,
		}
		assert.Equal(t, tt.link, p.ReLink(dlr), tt.name)
	}
}
