/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 */

package storage

import (
	"context"
)

type (
	// Entity defines entity context
	Entity interface {
		// ToString convert entity to string
		ToString() string

		// Key gets the key of the entity
		Key() string
	}

	// Relation defines relation context to generate the link
	// A -> B. The relation is defined in B. When it's create B entity
	// is linked with her parent depending the relation information that's provide B
	Relation interface {
		// ParentKey gets the key of the parent entity that contains this child entity
		ParentKey() string

		// Name gets the relation Name
		RelName() string

		// Link gets the relation entity that joins the parent and child. This link is used when is created the relation
		Link() Entity

		// LinkName gets the name which differentiates to the link. This name is used when is created
		// the doubly linked inverse relation
		LinkName() string

		// ReLink regenerates the link when change the relation name. When is called ReLink this entity
		// has the new information to change. 'dlr' parameter provides the relation information with old values
		ReLink(dlr DLRel) Entity

		// Empty builds empty relation entity
		Empty() EntityRelation
	}

	// EntityRelation groups Entity and Relation
	EntityRelation interface {
		Entity
		Relation
	}
)

type (
	// CRUD defines the CRUD operations
	CRUD interface {
		// Put creates or updates the entity depending of transaction. This context is added to the transaction.
		// The relation key cannot be changed
		// See Txn interface
		Put(entity Entity) (OpeWrap, error)

		// PutRaw puts the key and value depending of transaction. This context is added to the transaction.
		PutRaw(key string, value string) OpeWrap

		// Remove deletes the entity by key. This context is added to the transaction.
		// See Txn interface
		Remove(key string) OpeWrap

		// Get gets the entity into entity param
		Get(ctx context.Context, key string, entity Entity) (bool, error)

		// Exists if the key exists return true
		Exists(ctx context.Context, key string) (bool, error)

		// StartKey lists all entities that start by key with the limit of the top parameter.
		// Top = 0 is configured as unlimited
		StartKey(ctx context.Context, key string, top int, empty func() Entity) ([]Entity, error)

		// Range lists all entities that is greater than skey and ended by eKey with the limit of the top parameter.
		// Top = 0 is configured as unlimited
		Range(ctx context.Context, skey string, ekey string, top int, empty func() Entity) ([]Entity, error)

		// RangeRaw lists all keys and values that is greater than skey and ended by eKey with the limit of the top parameter.
		RangeRaw(ctx context.Context, skey string, ekey string, top int) (map[string]string, error)

		// Close closes resources
		Close() error
	}

	// Txn defines the transaction operations
	Txn interface {
		// Find checks if exists the keyValue. If it is found does DoFound or else does DoNotFound into commit
		Find(keyValue string)

		// DoFound saves the operations to transaction if it is found into commit
		DoFound(ope OpeWrap)

		// DoNotFound saves the operations to transaction if it is not found into commit
		DoNotFound(ope OpeWrap)

		// Commit commits the transaction. If it is returned true the transaction is successfully
		Commit(ctx context.Context) (bool, error)
	}
)

// Integration defines the functions to test
type Integration interface {
	// Store gets a store from integration
	Store() CRUD

	// Close closes the connections
	Close()
}
