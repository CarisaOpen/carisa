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

package relation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstSpace_ToString(t *testing.T) {
	i := InstSpace{
		ID:      "key",
		Name:    "name",
		SpaceID: "1",
	}
	assert.Equal(t, "inst-space-link: ID:key, Name:name", i.ToString())
}

func TestInstSpace_Key(t *testing.T) {
	i := InstSpace{
		ID:      "key",
		SpaceID: "1",
	}
	assert.Equal(t, i.ID, i.Key())
}

func TestSpaceEnte_ToString(t *testing.T) {
	s := SpaceEnte{
		ID:     "key",
		Name:   "name",
		EnteID: "1",
	}
	assert.Equal(t, "space-ente-link: ID:key, Name:name", s.ToString())
}

func TestSpaceEnte_Key(t *testing.T) {
	i := SpaceEnte{
		ID:     "key",
		EnteID: "1",
	}
	assert.Equal(t, i.ID, i.Key())
}

func TestEnteEnteProp_ToString(t *testing.T) {
	s := EnteEnteProp{
		ID:         "key",
		Name:       "name",
		EntePropID: "1",
	}
	assert.Equal(t, "ente-enteprop-link: ID:key, Name:name", s.ToString())
}

func TestEnteEnteProp_Key(t *testing.T) {
	i := EnteEnteProp{
		ID:         "key",
		EntePropID: "1",
	}
	assert.Equal(t, i.ID, i.Key())
}

func TestSpaceCategory_ToString(t *testing.T) {
	s := SpaceCategory{
		ID:    "key",
		Name:  "name",
		CatID: "1",
	}
	assert.Equal(t, "space-category-link: ID:key, Name:name", s.ToString())
}

func TestSpaceCategory_Key(t *testing.T) {
	s := SpaceCategory{
		ID:    "key",
		CatID: "1",
	}
	assert.Equal(t, s.ID, s.Key())
}

func TestHierarchy_ToString(t *testing.T) {
	h := Hierarchy{
		ID:     "key",
		Name:   "name",
		LinkID: "1",
	}
	assert.Equal(t, "hierarchy-link: ID:key, Name:name", h.ToString())
}

func TestHierarchy_Key(t *testing.T) {
	h := Hierarchy{
		ID:     "key",
		LinkID: "1",
	}
	assert.Equal(t, h.ID, h.Key())
}

func TestCategoryProp_ToString(t *testing.T) {
	c := CategoryProp{
		ID:        "key",
		Name:      "name",
		CatPropID: "1",
	}
	assert.Equal(t, "category-catprop-link: ID:key, Name:name", c.ToString())
}

func TestCategoryProp_Key(t *testing.T) {
	c := Hierarchy{
		ID:     "key",
		LinkID: "1",
	}
	assert.Equal(t, c.ID, c.Key())
}
