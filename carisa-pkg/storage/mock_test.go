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

func TestErrMockCRUD_Activate(t *testing.T) {
	m := ErrMockCRUD{}
	m.Activate("Create", "Close")
	_, err := m.Create(nil)
	assert.Error(t, err, "Create")
	err = m.Close()
	assert.Error(t, err, "Close")
}

func TestErrMockCRUD_ActivateMethodNotFound(t *testing.T) {
	m := ErrMockCRUD{}
	assert.Panics(t, func() { m.Activate("create") })
}

func TestErrMockCRUD_Clear(t *testing.T) {
	m := ErrMockCRUD{}
	m.Clear()
	assert.False(t, m.create, "Create")
	assert.False(t, m.close, "Close")
}

func TestErrMockTxn_Activate(t *testing.T) {
	m := ErrMockTxn{}
	m.Activate("Commit")
	_, err := m.Commit(nil)
	assert.Error(t, err, "Create")
}

func TestErrMockTxn_ActivateMethodNotFound(t *testing.T) {
	m := ErrMockTxn{}
	assert.Panics(t, func() { m.Activate("create") })
}

func TestErrMockTxn_Clear(t *testing.T) {
	m := ErrMockTxn{}
	m.Clear()
	assert.False(t, m.commit, "Commit")
}
