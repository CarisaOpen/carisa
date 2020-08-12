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
	assert.Equal(t, "instSpaceLink: ID:key, Name:name", i.ToString())
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
	assert.Equal(t, "spaceEnteLink: ID:key, Name:name", s.ToString())
}

func TestSpaceEnte_Key(t *testing.T) {
	i := SpaceEnte{
		ID:     "key",
		EnteID: "1",
	}
	assert.Equal(t, i.ID, i.Key())
}
