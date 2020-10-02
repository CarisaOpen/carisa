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

package instance

import (
	"testing"

	"github.com/carisa/internal/api/entity"

	"github.com/carisa/pkg/strings"

	"github.com/stretchr/testify/assert"
)

func TestInstance_ToString(t *testing.T) {
	i := New()
	i.Name = "name"
	assert.Contains(t, strings.Concat("instance: ID:", i.Key(), ", Name:", i.Name), i.ToString())
}

func TestInstance_Key(t *testing.T) {
	i := New()
	assert.Equal(t, i.ID.String(), i.Key())
}

func TestInstnace_Nominative(t *testing.T) {
	i := Instance{}
	assert.Equal(t, entity.Descriptor{}, i.Nominative())
}