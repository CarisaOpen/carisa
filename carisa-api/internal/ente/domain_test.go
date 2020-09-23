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

	"github.com/carisa/api/internal/entity"

	"github.com/carisa/api/internal/relation"

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
	name := "name"

	tests := []struct {
		name  string
		ente  Ente
		catId string
		link  storage.Entity
	}{
		{
			name: "Space with ente",
			ente: Ente{
				Descriptor: entity.Descriptor{
					ID:   enteID,
					Name: "name",
				},
				SpaceID: parentID,
				catID:   "",
			},
			link: &relation.SpaceEnte{
				ID:     strings.Concat(parentID.String(), relation.SpaceEnteLn, name, enteID.String()),
				Name:   name,
				EnteID: enteID.String(),
			},
		},
		{
			name: "Category with ente",
			ente: Ente{
				Descriptor: entity.Descriptor{
					ID:   enteID,
					Name: "name",
				},
				catID: parentID.String(),
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

func TestCategory_LinkName(t *testing.T) {
	c := New()

	tests := []struct {
		name  string
		catID string
		typen string
	}{
		{
			name:  "Category -> Ente",
			typen: relation.SpaceEnteLn,
		},
		{
			name:  "Space -> Ente",
			catID: "1",
			typen: relation.CatEnteLn,
		},
	}

	for _, tt := range tests {
		c.catID = tt.catID
		assert.Equal(t, tt.typen, c.LinkName(), tt.name)
	}
}

func TestEnte_ReLink(t *testing.T) {
	c := New()
	parentID := xid.New().String()

	tests := []struct {
		name string
		tn   string
		link storage.Entity
	}{
		{
			name: "Space with ente",
			tn:   relation.SpaceEnteLn,
			link: &relation.SpaceEnte{
				ID:     strings.Concat(parentID, relation.SpaceEnteLn, c.Name, c.Key()),
				Name:   c.Name,
				EnteID: c.ID.String(),
			},
		},
		{
			name: "Category with ente",
			tn:   relation.CatEnteLn,
			link: &relation.Hierarchy{
				ID:       strings.Concat(parentID, c.Name, c.Key()),
				Name:     c.Name,
				Category: false,
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

func TestEnteEnteProp_Field(t *testing.T) {
	e := NewProp()
	assert.Equal(t, e.Type, entity.Integer)
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
	e := NewProp()
	e.EnteID = xid.New()

	link := relation.EnteEnteProp{
		ID:         strings.Concat(e.EnteID.String(), relation.EntePropLn, e.Name, e.Key()),
		Name:       e.Name,
		EntePropID: e.ID.String(),
	}

	assert.Equal(t, &link, e.Link())
}

func TestEnteEnteProp_LinkName(t *testing.T) {
	c := NewProp()
	assert.Equal(t, relation.EntePropLn, c.LinkName())
}

func TestEnteEnteProp_ReLink(t *testing.T) {
	c := NewProp()
	parentID := xid.New().String()

	link := relation.EnteEnteProp{
		ID:         strings.Concat(parentID, relation.EntePropLn, c.Name, c.Key()),
		Name:       c.Name,
		EntePropID: c.ID.String(),
	}

	dlr := storage.DLRel{
		ParentID: parentID,
		Type:     relation.EntePropLn,
	}
	assert.Equal(t, &link, c.ReLink(dlr))
}
