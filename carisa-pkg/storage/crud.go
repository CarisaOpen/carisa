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
	// Create creates the entity into of the store
	Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error)

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

// Put implements CrudOperation.Put
func (c *crudOperation) Create(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	create, err := c.store.Put(entity)
	if err != nil {
		c.log.ErrorE(err, loc)
		return false, err
	}

	txn.DoNotFound(create)

	ctx, cancel := storeTimeout()
	ok, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "commit creating", loc, logging.String(reflect.TypeOf(entity).Name(), entity.ToString()))
	}

	return ok, nil
}

// Update implements CrudOperation.Put
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
