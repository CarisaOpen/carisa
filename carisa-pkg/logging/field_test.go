/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	const key = "key"
	value := "value"

	f := String(key, value)

	assert.Equal(t, f.key, key)
	assert.Equal(t, f.tpy, stringType)
	assert.Equal(t, f.stringV, value)
}

func TestBool(t *testing.T) {
	const key = "key"

	f := Bool(key, true)

	assert.Equal(t, f.key, key)
	assert.Equal(t, f.tpy, boolType)
	assert.True(t, f.boolV)
}
