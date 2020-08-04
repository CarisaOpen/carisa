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
	// If the entity was created returns true in the first param returned.
	// If the parent exists into store returns true in the second param returned.
	CreateWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error)

	// Put creates if exist the entity or updates if not exists the entity into storage
	Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error)

	// PutWithRel creates or updates the entity and the relation entity that joins the parent and child
	// The relation is updated if the relation name has changed
	// If the entity was updated returns true in the first param returned otherwise it is created.
	// If the parent exists into store returns true in the second param returned.
	PutWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error)
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
	return c.create(loc, storeTimeout, entity, false)
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
	created, err := c.create(loc, storeTimeout, entity, true)
	return created, true, err
}

// create creates the entity and if the relation exists entity also is created
func (c *crudOperation) create(loc string, storeTimeout StoreWithTimeout, entity Entity, isRel bool) (bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	create, err := c.store.Put(entity)
	if err != nil {
		return false, c.log.ErrWrap(err, "creating", loc)
	}

	txn.DoNotFound(create)
	// If the entity is a relation is inserted in the same transaction
	if isRel {
		rel, _ := entity.(EntityRelation)
		link := rel.Link()
		create, err := c.store.Put(&link)
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
	updated, err := c.put(loc, storeTimeout, entity, false)
	return updated, err
}

// PutWithRel implements CrudOperation.PutWithRel
func (c *crudOperation) PutWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error) {
	found, err := c.existsParent(loc, storeTimeout, entity)
	if err != nil {
		return false, false, err
	}
	if !found {
		return false, false, nil
	}
	updated, err := c.put(loc, storeTimeout, entity, true)
	return updated, true, err
}

func (c *crudOperation) put(loc string, storeTimeout StoreWithTimeout, entity Entity, isRel bool) (bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	put, err := c.store.Put(entity)
	if err != nil {
		c.log.ErrorE(err, loc)
		return false, err
	}

	ctx, cancel := storeTimeout()
	var link Link
	var rel EntityRelation
	if isRel {
		rel, _ = entity.(EntityRelation)
		link = rel.Link()
	}

	// Update entity
	txn.DoFound(put)
	// If the relation is passed by param and the entity exists is updated in the same transaction
	if isRel {
		if err := c.updateRelation(ctx, loc, rel, entity, txn, link); err != nil {
			return false, err
		}
	}

	// Create entity
	txn.DoNotFound(put)
	// If the relation is passed by param is inserted in the same transaction
	if isRel {
		create, err := c.store.Put(&link)
		if err != nil {
			return false, c.log.ErrWrap(err, "creating relation", loc)
		}
		txn.DoNotFound(create)
	}

	updated, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "commit putting", loc, logging.String(reflect.TypeOf(entity).Name(), entity.ToString()))
	}

	return updated, nil
}

func (c *crudOperation) updateRelation(
	ctx context.Context,
	loc string,
	rel EntityRelation,
	entity Entity,
	txn Txn,
	link Link) error {
	oldEntity := rel.Empty()
	found, err := c.store.Get(ctx, entity.Key(), oldEntity)
	if err != nil {
		return c.log.ErrWrap(err, "getting the entity", loc)
	}
	// Even if the entity is not found later it will be created with with DoNotFound
	if found {
		name := rel.RelName()
		if len(name) != 0 && name != oldEntity.RelName() {
			// it removes old relation and creating new relation when change name. Name is part of key
			txn.DoFound(c.store.Remove(oldEntity.RelKey()))
			updRel, err := c.store.Put(&link)
			if err != nil {
				return c.log.ErrWrap(err, "updating relation", loc)
			}
			txn.DoFound(updRel)
		}
	}
	return nil
}

func (c *crudOperation) existsParent(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, error) {
	ctx, cancel := storeTimeout()
	found, err := c.store.Exists(ctx, entity.ParentKey())
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "finding parent", loc, logging.String("Parent key", entity.ParentKey()))
	}
	return found, nil
}
