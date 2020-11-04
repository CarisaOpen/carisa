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

package space

import (
	"testing"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/entity"

	"github.com/carisa/internal/api/relation"

	"github.com/rs/xid"

	"github.com/carisa/pkg/strings"

	"github.com/stretchr/testify/assert"
)

func TestSpace_ToString(t *testing.T) {
	s := New()
	assert.Equal(t, strings.Concat("space: ID:", s.Key(), ", name:", s.Name), s.ToString())
}

func TestSpace_Key(t *testing.T) {
	s := New()
	assert.Equal(t, entity.SpaceKey(s.ID), s.Key())
}

func TestSpace_Nominative(t *testing.T) {
	s := Space{}
	assert.Equal(t, entity.Descriptor{}, s.Nominative())
}

func TestSpace_ParentKey(t *testing.T) {
	s := New()
	s.InstID = xid.New()
	assert.Equal(t, entity.InstKey(s.InstID), s.ParentKey())
}

func TestSpace_Empty(t *testing.T) {
	s := New()
	assert.Equal(t, &Space{}, s.Empty())
}

func TestSpace_Link(t *testing.T) {
	s := New()
	s.InstID = xid.New()

	link := relation.InstSpace{
		ID:      strings.Concat(entity.InstKey(s.InstID), relation.InstSpaceLn, s.Name, s.Key()),
		Name:    s.Name,
		SpaceID: s.ID.String(),
	}

	assert.Equal(t, &link, s.Link())
}

func TestSpace_LinkName(t *testing.T) {
	c := New()
	assert.Equal(t, relation.InstSpaceLn, c.LinkName())
}

func TestSpace_ReLink(t *testing.T) {
	c := New()
	parentID := xid.New()

	link := relation.InstSpace{
		ID:      strings.Concat(entity.InstKey(parentID), relation.InstSpaceLn, c.Name, c.Key()),
		Name:    c.Name,
		SpaceID: c.ID.String(),
	}

	dlr := storage.DLRel{
		ParentID: entity.InstKey(parentID),
		Type:     relation.InstSpaceLn,
	}
	assert.Equal(t, &link, c.ReLink(dlr))
}
