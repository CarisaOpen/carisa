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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type Object struct {
	ID     string
	Name   string
	Value  int
	Parent string
}

func (o Object) ToString() string {
	return o.ID
}

func (o Object) Key() string {
	return o.ID
}

func (o Object) RelKey() string {
	return strings.Concat(o.ParentKey(), o.Name, o.Key())
}

func (o Object) ParentKey() string {
	return o.Parent
}

func (o Object) RelName() string {
	return o.Name
}

func (o Object) Link() Entity {
	return &Link{
		ID:   o.RelKey(),
		Name: o.Name,
		Rel:  o.ID,
	}
}

func (o Object) Empty() EntityRelation {
	return &Object{}
}

type Link struct {
	ID   string
	Name string
	Rel  string
}

func (l Link) ToString() string {
	return strings.Concat("link: ID:", l.Key(), ", Name:", l.Name)
}

func (l Link) Key() string {
	return l.ID
}

func TestCRUDOperation_Store(t *testing.T) {
	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	assert.NotNil(t, oper.Store())
}

func TestCRUDOperation_Create(t *testing.T) {
	e := entity()

	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	ok, err := oper.Create("loc", storeTimeout, e)
	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
		var entityr Object
		found, err := oper.Store().Get(context.TODO(), e.Key(), &entityr)
		if assert.NoError(t, err) {
			assert.True(t, found, "Entity found")
			assert.Equal(t, e, entityr, "Entity saved")
		}
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
			err:   "creating: put",
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

func TestCRUDOperation_CreateWithRel(t *testing.T) {
	tests := []struct {
		parent  bool
		created bool
	}{
		{
			parent:  false,
			created: false,
		},
		{
			parent:  true,
			created: true,
		},
	}

	e := entity()

	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	for _, tt := range tests {
		if tt.parent {
			_, err := oper.Create("loc", storeTimeout, &Object{
				ID:    e.ParentKey(),
				Value: 1,
			})
			if err != nil {
				assert.Error(t, err)
				continue
			}
		}
		ok, foundParent, err := oper.CreateWithRel("loc", storeTimeout, e)
		if assert.NoError(t, err) {
			assert.Equal(t, tt.parent, foundParent, "Finding parent")
			assert.Equal(t, tt.created, ok, "Created")
			if tt.parent {
				var entityr Object
				found, err := oper.Store().Get(context.TODO(), e.Key(), &entityr)
				if assert.NoError(t, err) {
					assert.True(t, found, "Entity found")
					assert.Equal(t, e, entityr, "Entity saved")
				}
				var link Link
				found, err = oper.Store().Get(context.TODO(), e.RelKey(), &link)
				if assert.NoError(t, err) {
					assert.True(t, found, "Link found")
					assert.Equal(t, e.Link(), &link, "Link saved")
				}
			}
		}
	}
}

func TestCRUDOperation_CreateWithRelError(t *testing.T) {
	e := entity()

	tests := []struct {
		name  string
		mockS func(*ErrMockCRUD)
		err   string
	}{
		{
			name:  "Error finding parent",
			mockS: func(s *ErrMockCRUD) { s.Activate("Exists") },
			err:   "finding parent. Parent key: parentKey: exists",
		},
	}

	oper, store, _ := newCRUDOperMock()

	for _, tt := range tests {
		if tt.mockS != nil {
			tt.mockS(store)
		}
		_, _, err := oper.CreateWithRel("loc", storeTimeout, e)
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
		if assert.NoError(t, err, "Put failed") {
			assert.Equal(t, updated, tt.updated, "Updated")
			var entityr Object
			found, err := storef.Store().Get(context.TODO(), tt.e.ID, &entityr)
			if assert.NoError(t, err, "Commit failed") {
				assert.True(t, found, "Get entity")
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
		mockT func(*ErrMockTxn)
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
			err:   "put",
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

func TestCRUDOperation_PutWithRelation(t *testing.T) {
	tests := []struct {
		parent  bool
		updated bool
		e       *Object
	}{
		{
			parent:  true,
			updated: false,
			e: &Object{
				ID:     "key",
				Name:   "Name",
				Value:  1,
				Parent: "parentKey",
			},
		},
		{
			parent:  true,
			updated: true,
			e: &Object{
				ID:     "key",
				Name:   "name1",
				Value:  2,
				Parent: "parentKey",
			},
		},
		{
			parent:  true,
			updated: true,
			e: &Object{
				ID:     "key",
				Name:   "name1",
				Value:  3,
				Parent: "parentKey",
			},
		},
		{
			parent:  false,
			updated: false,
			e: &Object{
				Parent: "parentKey1",
			},
		},
	}

	storef := NewEctdIntegra(t)
	defer storef.Close()

	for _, tt := range tests {
		oper := newCRUDOper(storef)
		if tt.parent {
			_, err := oper.Create("loc", storeTimeout, &Object{
				ID:    tt.e.ParentKey(),
				Value: 1,
			})
			if err != nil {
				assert.Error(t, err)
				continue
			}
		}
		updated, foundParent, err := oper.PutWithRel("loc", storeTimeout, tt.e)
		if err != nil {
			assert.Error(t, err, "Put failed")
			continue
		}
		assert.Equal(t, tt.parent, foundParent, "Finding parent")
		assert.Equal(t, updated, tt.updated, "Updated")
		if tt.parent {
			var entityr Object
			found, err := storef.Store().Get(context.TODO(), tt.e.ID, &entityr)
			if err != nil {
				assert.Error(t, err, "Commit failed")
				continue
			}
			assert.True(t, found, "Get entity")
			assert.Equal(t, tt.e, &entityr, "Entity saved")
			var link Link
			found, err = storef.Store().Get(context.TODO(), tt.e.RelKey(), &link)
			if assert.NoError(t, err, "Commit failed") {
				assert.True(t, found, "Get link")
				assert.Equal(t, tt.e.Link(), &link, "Link saved")
			}
		}
	}
	ok, err := storef.Store().Exists(context.TODO(), tests[0].e.RelKey())
	if assert.NoError(t, err, "Obsolete relation removed") {
		assert.False(t, ok, "Obsolete relation removed")
	}
}

func TestCRUDOperation_PutWithRelError(t *testing.T) {
	e := entity()

	tests := []struct {
		name  string
		mockS func(*ErrMockCRUD)
		mockT func(*ErrMockTxn)
		err   string
	}{
		{
			name:  "Error finding parent",
			mockS: func(s *ErrMockCRUD) { s.Activate("Exists") },
			err:   "finding parent. Parent key: parentKey: exists",
		},
		{
			name:  "Error finding entity to see change",
			mockS: func(s *ErrMockCRUD) { s.Activate("Get") },
			err:   "getting the entity: get",
		},
	}

	oper, store, _ := newCRUDOperMock()

	for _, tt := range tests {
		if tt.mockS != nil {
			tt.mockS(store)
		}
		_, _, err := oper.PutWithRel("loc", storeTimeout, e)
		if assert.Error(t, err, tt.name) {
			assert.Equal(t, tt.err, err.Error())
		}
	}
}

func entity() Object {
	return Object{
		ID:     "key",
		Value:  1,
		Parent: "parentKey",
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
