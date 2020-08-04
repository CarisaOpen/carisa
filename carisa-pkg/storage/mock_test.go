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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrMockCRUD_Activate(t *testing.T) {
	m := ErrMockCRUD{}
	m.Activate("Put", "Remove", "Get", "Exists", "Close")
	_, err := m.Put(nil)
	assert.Error(t, err, "Put")
	_ = m.Remove("")
	assert.Error(t, err, "Remove")
	_, err = m.Get(context.TODO(), "", nil)
	assert.Error(t, err, "Get")
	_, err = m.Exists(context.TODO(), "")
	assert.Error(t, err, "Exists")
	err = m.Close()
	assert.Error(t, err, "Close")
}

func TestErrMockCRUD_ActivateMethodNotFound(t *testing.T) {
	m := ErrMockCRUD{}
	assert.Panics(t, func() { m.Activate("op") })
}

func TestErrMockCRUD_Clear(t *testing.T) {
	m := ErrMockCRUD{}
	m.Clear()
	assert.False(t, m.put, "Put")
	assert.False(t, m.put, "Remove")
	assert.False(t, m.get, "Get")
	assert.False(t, m.exists, "Exists")
	assert.False(t, m.close, "Close")
}

func TestErrMockTxn_Activate(t *testing.T) {
	m := ErrMockTxn{}
	m.Activate("Commit")
	_, err := m.Commit(context.TODO())
	assert.Error(t, err, "Commit")
}

func TestErrMockTxn_ActivateMethodNotFound(t *testing.T) {
	m := ErrMockTxn{}
	assert.Panics(t, func() { m.Activate("op") })
}

func TestErrMockTxn_Clear(t *testing.T) {
	m := ErrMockTxn{}
	m.Clear()
	assert.False(t, m.commit, "Commit")
}

func TestErrMockOper_Activate(t *testing.T) {
	m := NewErrMockCRUDOper()
	m.Activate("Create", "Put", "CreateWithRel")
	_, err := m.Create("", nil, nil)
	assert.Error(t, err, "Create")
	_, err = m.Put("", nil, nil)
	assert.Error(t, err, "Put")
	_, _, err = m.CreateWithRel("", nil, nil)
	assert.Error(t, err, "CreateWithError")
}

func TestErrMockOper_ActivateMethodNotFound(t *testing.T) {
	m := NewErrMockCRUDOper()
	assert.Panics(t, func() { m.Activate("op") })
}

func TestErrMockOper_Clear(t *testing.T) {
	m := NewErrMockCRUDOper()
	m.Clear()
	assert.False(t, m.create, "Create")
	assert.False(t, m.create, "Put")
	assert.False(t, m.create, "CreateWithError")
}
