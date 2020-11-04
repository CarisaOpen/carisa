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

package entity

import (
	"testing"

	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestInstKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchInstance, id.String()), InstKey(id))
}

func TestSpaceKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchSpace, id.String()), SpaceKey(id))
}

func TestEnteKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchEnte, id.String()), EnteKey(id))
}

func TestEntePropKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchEnteProp, id.String()), EntePropKey(id))
}

func TestCategoryKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchCategory, id.String()), CategoryKey(id))
}

func TestCategoryPropKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchCatProp, id.String()), CatPropKey(id))
}

func TestPluginKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchPlugin, id.String()), PluginKey(id))
}

func TestObjectKey(t *testing.T) {
	id := xid.New()
	assert.Equal(t, strings.Concat(SchObject, id.String()), ObjectKey(id))
}
