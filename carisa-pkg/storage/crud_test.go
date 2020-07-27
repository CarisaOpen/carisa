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
	"time"

	"github.com/carisa/pkg/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type Object struct {
	ID    string `json:"id,omitempty"`
	Value int    `json:"value,omitempty"`
}

func (o Object) ToString() string {
	return o.ID
}

func (o Object) Key() string {
	return o.ID
}

func TestCRUDOperation_Store(t *testing.T) {
	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	assert.NotNil(t, oper.Store())
}

func TestCRUDOperation_Create(t *testing.T) {
	o := entity()

	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	ok, err := oper.Create("loc", storeTimeout, o)

	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
	}
}

func TestCRUDOperation_CreateError(t *testing.T) {
	e := entity()

	tests := []struct {
		name  string
		mockS func(*ErrMockCRUD)
		mockT func(txn *ErrMockTxn)
		err   string
	}{
		{
			name:  "Error creating or putting",
			mockS: func(s *ErrMockCRUD) { s.Activate("Put") },
			err:   "create",
		},
		{
			name:  "Error commits transactions",
			mockS: func(s *ErrMockCRUD) { s.Clear() },
			mockT: func(s *ErrMockTxn) { s.Activate("Commit") },
			err:   "commit creating. Object: key: commit",
		},
	}

	oper, store, txn := newCRUDOperMock()

	for _, tt := range tests {
		if tt.mockS != nil {
			tt.mockS(store)
		}
		if tt.mockT != nil {
			tt.mockT(txn)
		}
		_, err := oper.Create("loc", storeTimeout, e)
		if assert.Error(t, err, tt.name) {
			assert.Equal(t, tt.err, err.Error())
		}
	}
}

func TestCRUDOperation_Put(t *testing.T) {
	tests := []struct {
		e       *Object
		updated bool
	}{
		{
			e: &Object{
				ID:    "key",
				Value: 1,
			},
			updated: false,
		},
		{
			e: &Object{
				ID:    "key",
				Value: 2,
			},
			updated: true,
		},
	}

	storef := NewEctdIntegra(t)
	defer storef.Close()

	for _, tt := range tests {
		oper := newCRUDOper(storef)
		updated, err := oper.Put("loc", storeTimeout, tt.e)
		if assert.NoErrorf(t, err, "Put failed") {
			assert.Equal(t, updated, tt.updated, "Updated")
			var entityr Object
			found, err := storef.Store().Get(context.TODO(), tt.e.ID, &entityr)
			assert.True(t, found, "Get entity")
			if assert.NoErrorf(t, err, "Commit failed") {
				assert.Equal(t, tt.e, &entityr, "Entity saved")
			}
		}
	}
}

func TestCRUDOperation_PutError(t *testing.T) {
	e := entity()

	tests := []struct {
		name  string
		mockS func(*ErrMockCRUD)
		mockT func(txn *ErrMockTxn)
		err   string
	}{
		{
			name:  "Error commits transactions",
			mockS: func(s *ErrMockCRUD) { s.Clear() },
			mockT: func(s *ErrMockTxn) { s.Activate("Commit") },
			err:   "commit putting. Object: key: commit",
		},
		{
			name:  "Error putting",
			mockS: func(s *ErrMockCRUD) { s.Activate("Put") },
			err:   "create",
		},
	}

	oper, store, txn := newCRUDOperMock()

	for _, tt := range tests {
		if tt.mockS != nil {
			tt.mockS(store)
		}
		if tt.mockT != nil {
			tt.mockT(txn)
		}
		_, err := oper.Put("loc", storeTimeout, e)
		if assert.Error(t, err, tt.name) {
			assert.Equal(t, tt.err, err.Error())
		}
	}
}

func entity() Object {
	return Object{
		ID:    "key",
		Value: 1,
	}
}

func storeTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func newCRUDOper(storef Integration) CrudOperation {
	core, _ := observer.New(zap.DebugLevel)
	log := logging.NewZapWrap(zap.New(core), logging.DebugLevel, "")

	return &crudOperation{
		store:    storef.Store(),
		log:      log,
		buildTxn: NewTxn,
	}
}

func newCRUDOperMock() (CrudOperation, *ErrMockCRUD, *ErrMockTxn) {
	core, _ := observer.New(zap.DebugLevel)
	log := logging.NewZapWrap(zap.New(core), logging.DebugLevel, "")

	store := &ErrMockCRUD{}
	txn := &ErrMockTxn{}

	return &crudOperation{
		store:    store,
		log:      log,
		buildTxn: func(s CRUD) Txn { return txn },
	}, store, txn
}
