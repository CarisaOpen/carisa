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

	"github.com/carisa/api/internal/relation"

	"github.com/rs/xid"

	"github.com/carisa/pkg/strings"

	"github.com/stretchr/testify/assert"
)

func TestSpace_ToString(t *testing.T) {
	s := NewSpace()
	assert.Equal(t, strings.Concat("space: ID:", s.Key(), ", name:", s.Name), s.ToString())
}

func TestSpace_Key(t *testing.T) {
	s := NewSpace()
	assert.Equal(t, s.ID.String(), s.Key())
}

func TestSpace_RelKey(t *testing.T) {
	s := Space{}
	s.Name = "name"
	assert.Equal(t, "00000000000000000000name00000000000000000000", s.RelKey())
}

func TestSpace_ParentKey(t *testing.T) {
	s := NewSpace()
	s.InstID = xid.New()
	assert.Equal(t, s.InstID.String(), s.ParentKey())
}

func TestSpace_Empty(t *testing.T) {
	s := NewSpace()
	assert.Equal(t, &Space{}, s.Empty())
}

func TestSpace_Link(t *testing.T) {
	s := NewSpace()
	s.InstID = xid.New()

	link := relation.InstSpace{
		ID:      strings.Concat(s.InstID.String(), s.Name, s.Key()),
		Name:    s.Name,
		SpaceID: s.ID.String(),
	}

	assert.Equal(t, &link, s.Link())
}
