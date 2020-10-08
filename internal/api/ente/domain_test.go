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
	"testing"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/entity"

	"github.com/carisa/internal/api/relation"

	"github.com/rs/xid"

	"github.com/carisa/pkg/strings"

	"github.com/stretchr/testify/assert"
)

func TestEnte_ToString(t *testing.T) {
	e := New()
	assert.Equal(t, strings.Concat("ente: ID:", e.Key(), ", name:", e.Name), e.ToString())
}

func TestEnte_Key(t *testing.T) {
	e := New()
	assert.Equal(t, e.ID.String(), e.Key())
}

func TestEnte_Nominative(t *testing.T) {
	e := Ente{}
	assert.Equal(t, entity.Descriptor{}, e.Nominative())
}

func TestEnte_ParentKey(t *testing.T) {
	e := New()
	e.SpaceID = xid.New()
	assert.Equal(t, e.SpaceID.String(), e.ParentKey())
}

func TestEnte_Empty(t *testing.T) {
	e := New()
	assert.Equal(t, &Ente{}, e.Empty())
}

func TestEnte_Link(t *testing.T) {
	enteID := xid.New()
	parentID := xid.New()
	name := "namele"

	tests := []struct {
		name  string
		ente  Ente
		catID string
		link  storage.Entity
	}{
		{
			name: "Space with prop",
			ente: Ente{
				Descriptor: entity.Descriptor{
					ID:   enteID,
					Name: "namele",
				},
				SpaceID: parentID,
				CatID:   "",
			},
			link: &relation.SpaceEnte{
				ID:     strings.Concat(parentID.String(), relation.SpaceEnteLn, name, enteID.String()),
				Name:   name,
				EnteID: enteID.String(),
			},
		},
		{
			name: "Category with prop",
			ente: Ente{
				Descriptor: entity.Descriptor{
					ID:   enteID,
					Name: "namele",
				},
				CatID: parentID.String(),
			},
			link: &relation.Hierarchy{
				ID:       strings.Concat(parentID.String(), name, enteID.String()),
				Name:     name,
				LinkID:   enteID.String(),
				Category: false,
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.link, tt.ente.Link(), tt.name)
	}
}

func TestEnte_LinkName(t *testing.T) {
	e := New()

	tests := []struct {
		name  string
		catID string
		typen string
	}{
		{
			name:  "Space -> Ente",
			typen: relation.SpaceEnteLn,
		},
		{
			name:  "Category -> Ente",
			catID: "1",
			typen: relation.CatEnteLn,
		},
	}

	for _, tt := range tests {
		e.CatID = tt.catID
		assert.Equal(t, tt.typen, e.LinkName(), tt.name)
	}
}

func TestEnte_ReLink(t *testing.T) {
	e := New()
	parentID := xid.New().String()

	tests := []struct {
		name string
		tn   string
		link storage.Entity
	}{
		{
			name: "Space -> Ente",
			tn:   relation.SpaceEnteLn,
			link: &relation.SpaceEnte{
				ID:     strings.Concat(parentID, relation.SpaceEnteLn, e.Name, e.Key()),
				Name:   e.Name,
				EnteID: e.ID.String(),
			},
		},
		{
			name: "Category -> Ente",
			tn:   relation.CatEnteLn,
			link: &relation.Hierarchy{
				ID:       strings.Concat(parentID, e.Name, e.Key()),
				Name:     e.Name,
				Category: false,
				LinkID:   e.ID.String(),
			},
		},
	}

	for _, tt := range tests {
		dlr := storage.DLRel{
			ParentID: parentID,
			Type:     tt.tn,
		}
		assert.Equal(t, tt.link, e.ReLink(dlr), tt.name)
	}
}

func TestEnteEnteProp_Field(t *testing.T) {
	e := NewProp()
	assert.Equal(t, e.GetType(), entity.Integer)
}

func TestEnteEnteProp_ToString(t *testing.T) {
	e := NewProp()
	assert.Equal(t, strings.Concat("ente-property: ID:", e.Key(), ", name:", e.Name), e.ToString())
}

func TestEnteEnteProp_Key(t *testing.T) {
	e := NewProp()
	assert.Equal(t, e.ID.String(), e.Key())
}

func TestEnteEnteProp_Nominative(t *testing.T) {
	e := Prop{}
	assert.Equal(t, entity.Descriptor{}, e.Nominative())
}

func TestEnteEnteProp_ParentKey(t *testing.T) {
	e := NewProp()
	e.EnteID = xid.New()
	assert.Equal(t, e.EnteID.String(), e.ParentKey())
}

func TestEnteEnteProp_Empty(t *testing.T) {
	e := NewProp()
	assert.Equal(t, &Prop{}, e.Empty())
}

func TestEnteEnteProp_Link(t *testing.T) {
	propID := xid.New()
	parentID := xid.New()
	name := "namele"

	tests := []struct {
		name string
		prop Prop
		link storage.Entity
	}{
		{
			name: "Ente -> Property",
			prop: Prop{
				Descriptor: entity.Descriptor{
					ID:   propID,
					Name: "namele",
				},
				EnteID:    parentID,
				CatPropID: "",
			},
			link: &relation.EnteProp{
				ID:         strings.Concat(parentID.String(), relation.EntePropLn, name, propID.String()),
				Name:       name,
				EntePropID: propID.String(),
			},
		},
		{
			name: "Category property -> Property",
			prop: Prop{
				Descriptor: entity.Descriptor{
					ID:   propID,
					Name: "namele",
				},
				CatPropID: parentID.String(),
			},
			link: &relation.CatPropProp{
				ID:       strings.Concat(parentID.String(), name, propID.String()),
				Name:     name,
				PropID:   propID.String(),
				Category: false,
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.link, tt.prop.Link(), tt.name)
	}
}

func TestEnteEnteProp_LinkName(t *testing.T) {
	p := NewProp()

	tests := []struct {
		name      string
		catPropID string
		typen     string
	}{
		{
			name:  "Ente -> Property",
			typen: relation.EntePropLn,
		},
		{
			name:      "Category property -> Property",
			catPropID: "1",
			typen:     relation.CatPropPropLn,
		},
	}

	for _, tt := range tests {
		p.CatPropID = tt.catPropID
		assert.Equal(t, tt.typen, p.LinkName(), tt.name)
	}
}

func TestEnteEnteProp_ReLink(t *testing.T) {
	p := NewProp()
	parentID := xid.New().String()

	tests := []struct {
		name string
		tn   string
		link storage.Entity
	}{
		{
			name: "Ente -> Property",
			tn:   relation.EntePropLn,
			link: &relation.EnteProp{
				ID:         strings.Concat(parentID, relation.EntePropLn, p.Name, p.Key()),
				Name:       p.Name,
				EntePropID: p.ID.String(),
			},
		},
		{
			name: "Category property -> Property",
			tn:   relation.CatPropPropLn,
			link: &relation.CatPropProp{
				ID:       strings.Concat(parentID, p.Name, p.Key()),
				Name:     p.Name,
				Category: false,
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
