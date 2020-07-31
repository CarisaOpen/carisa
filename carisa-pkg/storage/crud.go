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
	"reflect"

	"github.com/carisa/pkg/logging"
)

type StoreWithTimeout func() (context.Context, context.CancelFunc)

// NewTxn allows injection for test
type BuildTxn func(s CRUD) Txn

// CrudOperation defines the CRUD operations
type CrudOperation interface {
	// Store gets store
	Store() CRUD

	// Create creates the entity into of the store
	Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error)

	// CreateWithRel creates the entity and the relation entity that joins the parent and child
	// The child is entity param. The parent info is into of the entity param.
	// If the entity was created returns true in the first param returned.
	// If the parent exists into store returns true in the second param returned.
	CreateWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error)

	// Put creates if exist the entity or updates if not exists the entity into storage
	Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error)
}

// crudOperation defines the CRUD operations
type crudOperation struct {
	store    CRUD
	log      logging.Logger
	buildTxn BuildTxn
}

// NewCrudOperation builds the Crud operations
func NewCrudOperation(store CRUD, log logging.Logger, buildTxn BuildTxn) CrudOperation {
	return &crudOperation{
		store:    store,
		log:      log,
		buildTxn: buildTxn,
	}
}

// Store implements CrudOperation.Store
func (c *crudOperation) Store() CRUD {
	return c.store
}

// Create implements CrudOperation.Put
func (c *crudOperation) Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	return c.create(loc, storeTimeout, entity, nil)
}

// CreateWithRel implements CrudOperation.CreateWithRel
func (c *crudOperation) CreateWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error) {
	found, err := c.existsParent(loc, storeTimeout, entity)
	if err != nil {
		return false, false, err
	}
	if !found {
		return false, false, nil
	}
	created, err := c.create(loc, storeTimeout, entity, entity.Link())
	return created, true, err
}

// create creates the entity and if the relation exists entity also is created
func (c *crudOperation) create(loc string, storeTimeout StoreWithTimeout, entity Entity, rel *Link) (bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	create, err := c.store.Put(entity)
	if err != nil {
		return false, c.log.ErrWrap(err, "creating", loc)
	}

	txn.DoNotFound(create)

	// If the relation exists is inserted in the same transaction
	if rel != nil {
		create, err := c.store.Put(rel)
		if err != nil {
			return false, c.log.ErrWrap(err, "creating relation", loc)
		}
		txn.DoNotFound(create)
	}

	ctx, cancel := storeTimeout()
	ok, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "commit creating", loc, logging.String(reflect.TypeOf(entity).Name(), entity.ToString()))
	}

	return ok, nil
}

// Put implements CrudOperation.Put
func (c *crudOperation) Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	put, err := c.store.Put(entity)
	if err != nil {
		c.log.ErrorE(err, loc)
		return false, err
	}

	txn.DoFound(put)    // Update
	txn.DoNotFound(put) // Create

	ctx, cancel := storeTimeout()
	updated, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "commit putting", loc, logging.String(reflect.TypeOf(entity).Name(), entity.ToString()))
	}

	return updated, nil
}

func (c *crudOperation) existsParent(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, error) {
	ctx, cancel := storeTimeout()
	found, err := c.store.Exists(ctx, entity.ParentKey())
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "Finding parent", loc, logging.String("Parent key", entity.ParentKey()))
	}
	return found, nil
}
