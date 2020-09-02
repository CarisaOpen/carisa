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

func TestEnte_RelKey(t *testing.T) {
	e := Ente{}
	e.Name = "name"
	assert.Equal(t, "00000000000000000000name00000000000000000000", e.RelKey())
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
	e := New()
	e.SpaceID = xid.New()

	link := relation.SpaceEnte{
		ID:     strings.Concat(e.SpaceID.String(), e.Name, e.Key()),
		Name:   e.Name,
		EnteID: e.ID.String(),
	}

	assert.Equal(t, &link, e.Link())
}

func TestEnteEnteProp_Field(t *testing.T) {
	e := NewProp()
	assert.Equal(t, e.Type, Integer)
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
	e := EnteProp{}
	assert.Equal(t, entity.Descriptor{}, e.Nominative())
}

func TestEnteEnteProp_RelKey(t *testing.T) {
	e := EnteProp{}
	e.Name = "name"
	assert.Equal(t, "00000000000000000000name00000000000000000000", e.RelKey())
}

func TestEnteEnteProp_ParentKey(t *testing.T) {
	e := NewProp()
	e.EnteID = xid.New()
	assert.Equal(t, e.EnteID.String(), e.ParentKey())
}

func TestEnteEnteProp_Empty(t *testing.T) {
	e := NewProp()
	assert.Equal(t, &EnteProp{}, e.Empty())
}

func TestEnteEnteProp_Link(t *testing.T) {
	e := NewProp()
	e.EnteID = xid.New()

	link := relation.EnteEnteProp{
		ID:         strings.Concat(e.EnteID.String(), e.Name, e.Key()),
		Name:       e.Name,
		EntePropID: e.ID.String(),
	}

	assert.Equal(t, &link, e.Link())
}
