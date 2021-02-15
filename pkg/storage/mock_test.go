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
	m.Activate("Put", "PutRaw", "Remove", "Get", "GetRaw", "Exists", "StartKey", "Range", "RangeRaw", "Close")
	_, err := m.Put(nil)
	assert.Error(t, err, "Put")
	opew := m.PutRaw("", "")
	assert.Equal(t, opew, OpeWrap{}, "PutRaw")
	_ = m.Remove("")
	assert.Error(t, err, "Remove")
	_, err = m.Get(context.TODO(), "", nil)
	assert.Error(t, err, "Get")
	_, _, err = m.GetRaw(context.TODO(), "")
	assert.Error(t, err, "GetRaw")
	_, err = m.Exists(context.TODO(), "")
	assert.Error(t, err, "Exists")
	_, err = m.StartKey(context.TODO(), "", 0, nil)
	assert.Error(t, err, "StartKey")
	_, err = m.Range(context.TODO(), "", "", 0, nil)
	assert.Error(t, err, "Range")
	_, err = m.RangeRaw(context.TODO(), "", "", 0)
	assert.Error(t, err, "RangeRaw")
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
	assert.False(t, m.putRaw, "Put")
	assert.False(t, m.remove, "Remove")
	assert.False(t, m.get, "Get")
	assert.False(t, m.getRaw, "GetRaw")
	assert.False(t, m.exists, "Exists")
	assert.False(t, m.startKey, "StartKey")
	assert.False(t, m.rang, "Range")
	assert.False(t, m.rangRaw, "RangeRaw")
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
	m.Activate("Create", "Put", "CreateWithRel", "PutWithRel", "Update", "LinkTo", "ListDLR")
	_, err := m.Create("", nil, nil)
	assert.Error(t, err, "Create")
	_, err = m.Put("", nil, nil)
	assert.Error(t, err, "Put")
	_, _, err = m.CreateWithRel("", nil, nil)
	assert.Error(t, err, "CreateWithError")
	_, _, err = m.PutWithRel("", nil, nil)
	assert.Error(t, err, "PutWithRel")
	_, err = m.Update("", nil, nil, nil)
	assert.Error(t, err, "Update")
	_, _, _, err = m.LinkTo("", nil, nil, nil, "", nil)
	assert.Error(t, err, "LinkTo")
	_, err = m.ListDLR(nil, "")
	assert.Error(t, err, "ListDLR")
}

func TestErrMockOper_ActivateMethodNotFound(t *testing.T) {
	m := NewErrMockCRUDOper()
	assert.Panics(t, func() { m.Activate("op") })
}

func TestErrMockOper_Clear(t *testing.T) {
	m := NewErrMockCRUDOper()
	m.Clear()
	assert.False(t, m.create, "Create")
	assert.False(t, m.put, "Put")
	assert.False(t, m.createWithRel, "CreateWithError")
	assert.False(t, m.putWithRel, "PutWithError")
	assert.False(t, m.update, "Update")
	assert.False(t, m.connectTo, "LinkTo")
	assert.False(t, m.listDLR, "ListDLR")
}
