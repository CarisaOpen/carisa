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

func (o Object) ParentKey() string {
	return o.Parent
}

func (o Object) RelName() string {
	return o.Name
}

func (o Object) Link() Entity {
	return &Link{
		ID:   strings.Concat(o.ParentKey(), o.Name, o.Key()),
		Name: o.Name,
		Rel:  o.ID,
	}
}

func (o Object) LinkName() string {
	return "linkname"
}

func (o Object) ReLink(dlr DLRel) Entity {
	return o.Link()
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
	const parentKey = "parentKey"
	const name = "namec"

	tests := []struct {
		name    string
		e       Object
		parent  bool
		created bool
	}{
		{
			name: "Creating. Parent not found",
			e: Object{
				ID:     "key",
				Name:   name,
				Parent: "parentKey1",
			},
			parent:  false,
			created: false,
		},
		{
			name: "Creating.",
			e: Object{
				ID:     "key1",
				Name:   name,
				Value:  1,
				Parent: parentKey,
			},
			parent:  true,
			created: true,
		},
		{
			name: "Creating with virtual parent.",
			e: Object{
				ID:     "key2",
				Name:   name,
				Value:  2,
				Parent: Virtual,
			},
			parent:  true,
			created: true,
		},
	}

	storef := NewEctdIntegra(t)
	oper := newCRUDOper(storef)
	defer storef.Close()

	_, err := oper.Create("loc", storeTimeout, &Object{
		ID:    parentKey,
		Value: 1,
	})
	if err != nil {
		assert.Error(t, err)
		return
	}

	for _, tt := range tests {
		ok, foundParent, err := oper.CreateWithRel("loc", storeTimeout, &tt.e)
		if err != nil {
			assert.Error(t, err)
			continue
		}

		assert.Equal(t, tt.parent, foundParent, strings.Concat(tt.name, "Finding parent"))
		assert.Equal(t, tt.created, ok, strings.Concat(tt.name, "Created"))

		if tt.parent {
			var entityr Object
			found, err := oper.Store().Get(context.TODO(), tt.e.Key(), &entityr)
			if assert.NoError(t, err) {
				assert.True(t, found, strings.Concat(tt.name, "Entity found"))
				assert.Equal(t, tt.e, entityr, strings.Concat(tt.name, "Entity saved"))
			}

			lnk := tt.e.Link()
			var link Link
			found, err = oper.Store().Get(context.TODO(), lnk.Key(), &link)
			if assert.NoError(t, err) {
				assert.True(t, found, strings.Concat(tt.name, "Link found"))
				assert.Equal(t, lnk, &link, strings.Concat(tt.name, "Link saved"))
			}

			var dlr DLRel
			found, err = oper.Store().Get(context.TODO(), strings.Concat(entityr.ID, dlrSep, entityr.Parent), &dlr)
			if assert.NoError(t, err) {
				assert.True(t, found, strings.Concat(tt.name, "DlR found"))
				assert.Equal(t, DLRel{
					ChildID:  entityr.ID,
					ParentID: entityr.Parent,
					Type:     entityr.LinkName(),
					Pointer:  link.Key(),
				}, dlr, strings.Concat(tt.name, "DLR saved"))
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
		_, _, err := oper.CreateWithRel("loc", storeTimeout, &e)
		if assert.Error(t, err, tt.name) {
			assert.Equal(t, tt.err, err.Error())
		}
	}
}

func TestCRUDOperation_Put(t *testing.T) {
	tests := []struct {
		name    string
		e       *Object
		updated bool
	}{
		{
			name: "Creating.",
			e: &Object{
				ID:    "key",
				Value: 1,
			},
			updated: false,
		},
		{
			name: "Updating.",
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
		if assert.NoError(t, err, strings.Concat(tt.name, "Put failed")) {
			assert.Equal(t, updated, tt.updated, "Updated")

			var entityr Object
			found, err := storef.Store().Get(context.TODO(), tt.e.ID, &entityr)
			if assert.NoError(t, err, strings.Concat(tt.name, "Get entity")) {
				assert.True(t, found, strings.Concat(tt.name, "Get entity"))
				assert.Equal(t, tt.e, &entityr, strings.Concat(tt.name, "Entity saved"))
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

func TestCRUDOperation_Update(t *testing.T) {
	tests := []struct {
		name  string
		e     *Object
		found bool
	}{
		{
			name: "Not found.",
			e: &Object{
				ID:    "key1",
				Value: 1,
			},
			found: false,
		},
		{
			name: "Updated.",
			e: &Object{
				ID:    "key",
				Value: 5,
			},
			found: true,
		},
	}

	storef := NewEctdIntegra(t)
	defer storef.Close()

	o := Object{
		ID:    "key",
		Value: 2,
	}

	oper := newCRUDOper(storef)
	_, err := oper.Put("loc", storeTimeout, &o)
	if err != nil {
		assert.Error(t, err, "Creating entity")
		return
	}

	for _, tt := range tests {
		cpy := *tt.e
		found, err := oper.Update("loc", storeTimeout, tt.e, func(e Entity) {
			e.(*Object).Value = cpy.Value
		})
		if assert.NoError(t, err, strings.Concat(tt.name, "Updating")) {
			assert.Equal(t, tt.found, found, strings.Concat(tt.name, "Found"))
			if tt.found {
				var res Object
				_, err := storef.Store().Get(context.TODO(), tt.e.ID, &res)
				if assert.NoError(t, err, strings.Concat(tt.name, "Get entity")) {
					assert.Equal(t, cpy, res, strings.Concat(tt.name, "Entity updated"))
				}
			}
		}
	}
}

func TestCRUDOperation_UpdateError(t *testing.T) {
	e := entity()

	oper, store, _ := newCRUDOperMock()

	store.Activate("Get")

	_, err := oper.Update("loc", storeTimeout, e, func(entity Entity) {})
	if assert.Error(t, err, "Update") {
		assert.Equal(t, "get", err.Error())
	}
}

func TestCRUDOperation_PutWithRelation(t *testing.T) {
	tests := []struct {
		name    string
		parent  bool
		updated bool
		e       *Object
	}{
		{
			name:    "Creating.",
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
			name:    "Updating name.",
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
			name:    "Updating value. No replace relation.",
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
			name:    "Creating. Parent not found.",
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
			assert.Error(t, err, strings.Concat(tt.name, "Put failed"))
			continue
		}
		assert.Equal(t, tt.parent, foundParent, strings.Concat(tt.name, "Finding parent"))
		assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Updated"))
		if tt.parent {
			var entityr Object
			found, err := storef.Store().Get(context.TODO(), tt.e.ID, &entityr)
			if err != nil {
				assert.Error(t, err, strings.Concat(tt.name, "Error getting entity"))
				continue
			}
			assert.True(t, found, strings.Concat(tt.name, "Get entity"))
			assert.Equal(t, tt.e, &entityr, strings.Concat(tt.name, "Entity saved"))

			lnk := tt.e.Link()
			var link Link
			found, err = storef.Store().Get(context.TODO(), lnk.Key(), &link)
			if assert.NoError(t, err, strings.Concat(tt.name, "Error getting link")) {
				assert.True(t, found, strings.Concat(tt.name, "Get link"))
				assert.Equal(t, lnk, &link, strings.Concat(tt.name, "Link saved"))
			}

			var dlr DLRel
			found, err = oper.Store().Get(context.TODO(), strings.Concat(entityr.ID, dlrSep, entityr.Parent), &dlr)
			if assert.NoError(t, err, strings.Concat(tt.name, "Error getting DLR")) {
				assert.True(t, found, strings.Concat(tt.name, "DlR found"))
				assert.Equal(t, DLRel{
					ChildID:  entityr.ID,
					ParentID: entityr.Parent,
					Type:     entityr.LinkName(),
					Pointer:  link.Key(),
				}, dlr, strings.Concat(tt.name, "DLR saved"))
			}
		}
	}
	lnk := tests[0].e.Link()
	ok, err := storef.Store().Exists(context.TODO(), lnk.Key())
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
		_, _, err := oper.PutWithRel("loc", storeTimeout, &e)
		if assert.Error(t, err, tt.name) {
			assert.Equal(t, tt.err, err.Error())
		}
	}
}

func TestCRUDOperation_LinkTo(t *testing.T) {
	tests := []struct {
		name      string
		kchild    Object
		kparentID string
		pfound    bool
		cfound    bool
	}{
		{
			name:   "Parent not found.",
			kchild: Object{ID: "k"},
			pfound: false,
			cfound: false,
		},
		{
			name:      "Child not found.",
			kchild:    Object{ID: "key"},
			kparentID: "key1",
			pfound:    true,
			cfound:    false,
		},
		{
			name:      "Connected.",
			kchild:    Object{ID: "key"},
			kparentID: "key",
			pfound:    true,
			cfound:    true,
		},
	}

	storef := NewEctdIntegra(t)
	defer storef.Close()

	oper := newCRUDOper(storef)
	_, err := sampleConnectTo(t, oper)
	if err != nil {
		assert.Error(t, err, "Inserting sample")
	}

	for _, tt := range tests {
		sfound, tfound, link, err := oper.LinkTo("loc", storeTimeout, nil, &tt.kchild, tt.kparentID, func(child Entity) {
			child.(*Object).Parent = tt.kparentID
		})
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.pfound, sfound, "Source")
			assert.Equal(t, tt.cfound, tfound, "Target")
			if tt.cfound && tt.pfound {
				ttlink := tt.kchild.Link()
				assert.Equal(t, ttlink, link, "Relation")

				var rel Link
				_, err := storef.Store().Get(context.TODO(), ttlink.Key(), &rel)
				if assert.NoError(t, err, tt.name) {
					assert.Equal(t, ttlink, &rel, strings.Concat(tt.name, "Relation created"))
				}

				var dlr DLRel
				found, err := oper.Store().Get(context.TODO(), strings.Concat(tt.kchild.ID, dlrSep, tt.kparentID), &dlr)
				if assert.NoError(t, err, strings.Concat(tt.name, "Error getting DLR")) {
					assert.True(t, found, strings.Concat(tt.name, "DlR found"))
					assert.Equal(t, DLRel{
						ChildID:  tt.kchild.ID,
						ParentID: tt.kparentID,
						Type:     tt.kchild.LinkName(),
						Pointer:  link.Key(),
					}, dlr, strings.Concat(tt.name, "DLR saved"))
				}
			}
		}
	}
}

func TestCRUDOperation_LinkTo_Txn(t *testing.T) {
	storef := NewEctdIntegra(t)
	defer storef.Close()

	ot := &Object{
		ID:   "keyt",
		Name: "name",
	}

	oper := newCRUDOper(storef)
	o, err := sampleConnectTo(t, oper)
	if err != nil {
		assert.Error(t, err, "Inserting sample")
	}

	txn := NewTxn(oper.Store())

	put, err := oper.Store().Put(ot)
	if err != nil {
		assert.Error(t, err, "Put txn")
		return
	}
	txn.DoNotFound(put)

	_, _, _, err = oper.LinkTo("loc", storeTimeout, txn, o, o.ID, func(child Entity) {
		child.(*Object).Parent = o.ID
	})
	if assert.NoError(t, err, "LinkTo") {
		ctx, cancel := storeTimeout()
		_, err := txn.Commit(ctx)
		cancel()
		if err != nil {
			assert.Error(t, err, "Commit")
			return
		}
		found, err := storef.Store().Exists(context.TODO(), o.Link().Key())
		if assert.NoError(t, err, "Relation") {
			assert.True(t, found, "Exists rel")
			found, err = storef.Store().Exists(context.TODO(), ot.ID)
			if assert.NoError(t, err, "Transaction entity") {
				assert.True(t, found, "Exists transaction entity")
			}
		}
	}
}

func TestCRUDOperation_ListDLR(t *testing.T) {
	storef := NewEctdIntegra(t)
	defer storef.Close()
	oper := newCRUDOper(storef)

	const childID = "C1"
	dlrTest := []DLRel{
		{
			ChildID:  childID,
			ParentID: "P1",
			Type:     "T1",
			Pointer:  "P1",
		},
		{
			ChildID:  childID,
			ParentID: "P2",
			Type:     "T3",
			Pointer:  "P2",
		},
	}

	for i := range dlrTest {
		_, err := oper.Create("loc", storeTimeout, &dlrTest[i])
		if err != nil {
			assert.Error(t, err, "Creating DLR")
			return
		}
	}

	dlrs, err := oper.ListDLR(storeTimeout, childID)
	if assert.NoError(t, err, "Listing DLR") {
		for i, dlr := range dlrs {
			assert.Equal(t, &dlrTest[i], dlr)
		}
	}
}

func sampleConnectTo(t *testing.T, oper CrudOperation) (*Object, error) {
	o := &Object{
		ID:   "key",
		Name: "name",
	}

	_, err := oper.Put("loc", storeTimeout, o)
	if err != nil {
		assert.Error(t, err, "Put parent")
		return nil, err
	}

	return o, nil
}

func entity() Object {
	return Object{
		ID:     "key",
		Name:   "name",
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

	return NewCrudOperation(store, log, func(s CRUD) Txn { return txn }), store, txn
}
