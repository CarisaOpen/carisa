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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/logging"
)

const dlrSep = "#R#"

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
	// In addition to the relation it creates a doubled linked relation (dlr) between the child
	// the relation and the parent. This can useful to navigate in inverse order
	// P -> R -> C, where P: Parent entity, R: Relation, C: Child
	// If the entity was created returns true in the first param returned.
	// If the parent exists into store returns true in the second param returned.
	CreateWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error)

	// Put creates if exist the entity or updates if not exists the entity into storage
	Put(loc string, storeTimeout StoreWithTimeout, entity Entity) (bool, error)

	// Update updates the entity fields through of the 'upd' parameter
	// The entity is searched using entity.Key(). The entity parameter will be replaced with the entity stored in the DB.
	Update(loc string, storeTimeout StoreWithTimeout, entity Entity, upd func(entity Entity)) (bool, error)

	// PutWithRel creates or updates the entity and the relation entity that joins the parent and child
	// The relation is found if the relation 'name' has changed
	// In addition to the relation it creates a doubled linked relation (dlr) between the child
	// the relation and the parent. This can be useful to navigate in inverse order
	// P -> R -> C, where P: Parent entity, R: Relation, C: Child
	// If the entity was found returns true in the first param returned otherwise it is created.
	// If the parent exists into store returns true in the second param returned.
	PutWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error)

	// LinkTo creates relation that links parent and child
	// The child parameter must be filled the ID.
	// The 'fill' parameter is a function to complete the fields of the child. This function receive as parameter
	// the child entity of the DB
	// If the transaction is sent as parameter just can be configured in DoNotFound
	// If the child entity exist returns true in the first param returned otherwise return false.
	// If the parent entity exist returns true in the second param returned otherwise return false.
	LinkTo(
		loc string,
		storeTimeout StoreWithTimeout,
		txn Txn,
		child EntityRelation,
		parentID string,
		fill func(child Entity)) (bool, bool, Entity, error)

	// Return a DLR slice from child identifier
	ListDLR(storeTimeout StoreWithTimeout, childID string) ([]Entity, error)
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
		err := c.createRel(loc, txn, entity.(EntityRelation))
		if err != nil {
			return false, err
		}
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
	updated, _, err := c.put(loc, storeTimeout, entity, false)
	return updated, err
}

// PutWithRel implements CrudOperation.PutWithRel
func (c *crudOperation) PutWithRel(loc string, storeTimeout StoreWithTimeout, entity EntityRelation) (bool, bool, error) {
	return c.put(loc, storeTimeout, entity, true)
}

// Update implements CrudOperation.Update
func (c *crudOperation) Update(
	loc string,
	storeTimeout StoreWithTimeout,
	entity Entity,
	upd func(entity Entity)) (bool, error) {
	//
	ctx, cancel := storeTimeout()
	found, err := c.store.Get(ctx, entity.Key(), entity)
	cancel()
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	// Update the fields
	upd(entity)

	return c.Put(loc, storeTimeout, entity)
}

// ConnectTo implements CrudOperation.ConnectTo
func (c *crudOperation) LinkTo(
	loc string,
	storeTimeout StoreWithTimeout,
	txn Txn,
	child EntityRelation,
	parentID string,
	fill func(child Entity)) (bool, bool, Entity, error) {
	//
	ctx, cancel := storeTimeout()
	found, err := c.Store().Get(ctx, child.Key(), child)
	cancel()
	if err != nil {
		return false, false, nil, c.log.ErrWrap1(err, "finding child entity to link", loc, logging.String("key", child.Key()))
	}
	if !found {
		return false, false, nil, nil
	}

	found, err = c.exists(loc, storeTimeout, parentID)
	if err != nil {
		return false, false, nil, c.log.ErrWrap1(err, "finding parent entity to link", loc, logging.String("key", parentID))
	}
	if !found {
		return true, false, nil, nil
	}

	if txn == nil {
		txn = c.buildTxn(c.store)
	}

	fill(child)
	link := child.Link()

	txn.Find(link.Key())

	err = c.createRel(loc, txn, child)
	if err != nil {
		return true, true, nil, err
	}

	ctx, cancel = storeTimeout()
	_, err = txn.Commit(ctx)
	cancel()
	if err != nil {
		return true, true, nil, err
	}

	return true, true, link, nil
}

func (c *crudOperation) ListDLR(storeTimeout StoreWithTimeout, childID string) ([]Entity, error) {
	ctx, cancel := storeTimeout()
	es, err := c.store.Range(ctx, strings.Concat(childID, dlrSep), childID, 0, func() Entity {
		return &DLRel{}
	})
	cancel()
	return es, err
}

func (c *crudOperation) exists(loc string, storeTimeout StoreWithTimeout, id string) (bool, error) {
	ctx, cancel := storeTimeout()
	found, err := c.store.Exists(ctx, id)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap1(err, "finding entity", loc, logging.String("key", id))
	}
	return found, nil
}

func (c *crudOperation) put(loc string, storeTimeout StoreWithTimeout, entity Entity, isRel bool) (bool, bool, error) {
	txn := c.buildTxn(c.store)
	txn.Find(entity.Key())

	// If the relation is passed by param and the entity exists is found in the same transaction
	if isRel {
		found, err := c.updateRel(storeTimeout, loc, entity, txn)
		if err != nil {
			return false, false, err
		}
		if !found { // If the entity is new, checks the parent
			found, err = c.existsParent(loc, storeTimeout, entity.(EntityRelation))
			if err != nil {
				return false, false, err
			}
			if !found { // The parent must exist
				return false, false, nil
			}
		}
	}

	put, err := c.store.Put(entity)
	if err != nil {
		c.log.ErrorE(err, loc)
		return false, false, err
	}
	// Update entity
	txn.DoFound(put)

	// Create entity
	txn.DoNotFound(put)
	// If the relation is passed by param is inserted in the same transaction
	if isRel {
		err := c.createRel(loc, txn, entity.(EntityRelation))
		if err != nil {
			return false, false, err
		}
	}

	ctx, cancel := storeTimeout()
	updated, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, false,
			c.log.ErrWrap1(err, "commit putting", loc, logging.String(reflect.TypeOf(entity).Name(), entity.ToString()))
	}

	return updated, true, nil
}

// createRel creates the relation between the parent and child and adds a doubly linked relation from child to the relation
// and to the parent
func (c *crudOperation) createRel(loc string, txn Txn, relation EntityRelation) error {
	rel := relation.Link()

	putRel, err := c.store.Put(rel)
	if err != nil {
		return c.log.ErrWrap(err, "creating relation", loc)
	}
	txn.DoNotFound(putRel)

	idx := &DLRel{
		ChildID:  relation.Key(),
		ParentID: relation.ParentKey(),
		Type:     relation.LinkName(),
		Pointer:  rel.Key(),
	}
	putIdx, err := c.store.Put(idx)
	if err != nil {
		return c.log.ErrWrap(err, "creating doubly linked relation", loc)
	}
	txn.DoNotFound(putIdx)

	return nil
}

func (c *crudOperation) updateRel(storeTimeout StoreWithTimeout, loc string, entity Entity, txn Txn) (bool, error) {
	rel, _ := entity.(EntityRelation)
	oldEntity := rel.Empty()
	ctx, cancel := storeTimeout()
	found, err := c.store.Get(ctx, entity.Key(), oldEntity)
	cancel()
	if err != nil {
		return false, c.log.ErrWrap(err, "getting the entity", loc)
	}
	if !found {
		return false, nil
	}

	// Even if the entity is not found later it will be created with with DoNotFound
	name := rel.RelName()
	if len(name) != 0 && name != oldEntity.RelName() {
		// it removes old relation and creating new relation when change Name. Name is part of key
		ctx, cancel := storeTimeout()
		dlrs, err := c.store.Range(
			ctx,
			strings.Concat(entity.Key(), dlrSep),
			entity.Key(),
			0,
			func() Entity { return &DLRel{} })
		cancel()
		if err != nil {
			return false, c.log.ErrWrap(err, "finding dlr", loc)
		}

		// Iterates all doubly linked relation (dlr) to change the relations
		for _, dlre := range dlrs {
			dlr := dlre.(*DLRel)

			// Remove old relation
			txn.DoFound(c.store.Remove(dlr.Pointer))

			// Create the new relation with tne new name
			linkr := rel.ReLink(*dlr)
			updlink, err := c.store.Put(linkr)
			if err != nil {
				return true, c.log.ErrWrap(err, "updating changed relation", loc)
			}
			txn.DoFound(updlink)

			// Change the pointer to the new relation
			dlr.Pointer = linkr.Key()
			putDlr, err := c.store.Put(dlr)
			if err != nil {
				return true, c.log.ErrWrap(err, "updating doubly linked relation", loc)
			}
			txn.DoFound(putDlr)
		}
	}
	return true, nil
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

// DLRel is a doubly linked relation
// P -> R -> C, where P: Parent entity, R: Relation, C: Child
// C -> DLR -> R where C: Child, DLR: Doubly linked relation, R: Relation
// It allows navigate from child to parent
type DLRel struct {
	ChildID  string
	ParentID string
	Type     string
	Pointer  string
}

func (r *DLRel) ToString() string {
	return strings.Concat("Relation index. ChildID: (", r.ChildID, ", ParentID:", r.ParentID, ")")
}

func (r *DLRel) Key() string {
	return DLRKey(r.ChildID, r.ParentID)
}

// DLRKey gets DLR key
func DLRKey(childID string, parentID string) string {
	return strings.Concat(childID, dlrSep, parentID)
}
