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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLink_ToString(t *testing.T) {
	i := Link{
		ID:   "key",
		Name: "name",
		Rel:  "1",
	}
	assert.Equal(t, "link: ID:key, Name:name", i.ToString())
}

func TestLink_Key(t *testing.T) {
	i := Link{
		ID:  "key",
		Rel: "1",
	}
	assert.Equal(t, i.ID, i.Key())
}
