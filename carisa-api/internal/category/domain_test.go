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

	"github.com/carisa/api/internal/ente"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/api/internal/entity"

	"github.com/carisa/api/internal/relation"

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
	assert.Equal(t, c.ID.String(), c.Key())
}

func TestCategory_Nominative(t *testing.T) {
	c := Category{}
	assert.Equal(t, entity.Descriptor{}, c.Nominative())
}

func TestCategory_RelKey(t *testing.T) {
	c := Category{}
	c.Name = "namec"
	c.Root = true
	assert.Equal(t, "00000000000000000000Snamec00000000000000000000", c.RelKey())
}

func TestCategory_ParentKey(t *testing.T) {
	c := New()
	c.ParentID = xid.New()
	assert.Equal(t, c.ParentID.String(), c.ParentKey())
}

func TestCategory_SetParentKey(t *testing.T) {
	c := New()
	_ = c.SetParentKey(xid.New().String())
	assert.Equal(t, c.ParentID.String(), c.ParentKey())
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
				ID:    strings.Concat(c.ParentID.String(), "S", c.Name, c.Key()),
				Name:  c.Name,
				CatID: c.ID.String(),
			},
		},
		{
			name: "Category with others Category",
			root: false,
			link: &relation.Hierarchy{
				ID:       strings.Concat(c.ParentID.String(), "C", c.Name, c.Key()),
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

func TestCategoryCatProp_Field(t *testing.T) {
	c := NewProp()
	assert.Equal(t, c.Type, ente.None)
}

func TestCategoryCatProp_ToString(t *testing.T) {
	c := NewProp()
	assert.Equal(t, strings.Concat("category-property: ID:", c.Key(), ", name:", c.Name), c.ToString())
}

func TestCategoryCatProp_Key(t *testing.T) {
	c := NewProp()
	assert.Equal(t, c.ID.String(), c.Key())
}

func TestCategoryCatProp_Nominative(t *testing.T) {
	c := Prop{}
	assert.Equal(t, entity.Descriptor{}, c.Nominative())
}

func TestCategoryCatProp_RelKey(t *testing.T) {
	c := Prop{}
	c.Name = "namecp"
	assert.Equal(t, "00000000000000000000Pnamecp00000000000000000000", c.RelKey())
}

func TestCategoryCatProp_ParentKey(t *testing.T) {
	c := NewProp()
	c.CatID = xid.New()
	assert.Equal(t, c.CatID.String(), c.ParentKey())
}

func TestCategoryCatProp_SetParentKey(t *testing.T) {
	c := NewProp()
	_ = c.SetParentKey(xid.New().String())
	assert.Equal(t, c.CatID.String(), c.ParentKey())
}

func TestCategoryCatProp_Empty(t *testing.T) {
	c := NewProp()
	assert.Equal(t, &Prop{}, c.Empty())
}

func TestCategoryCatProp_Link(t *testing.T) {
	c := NewProp()
	c.CatID = xid.New()

	link := relation.CategoryProp{
		ID:        strings.Concat(c.CatID.String(), "P", c.Name, c.Key()),
		Name:      c.Name,
		CatPropID: c.ID.String(),
	}

	assert.Equal(t, &link, c.Link())
}
