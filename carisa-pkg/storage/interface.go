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

// Entity defines entity context
type (
	Entity interface {
		// ToString convert entity to string
		ToString() string

		// GetKey gets the key
		GetKey() string
	}
)

type (
	// CRUD defines the CRUD operations
	CRUD interface {
		// Create creates the context to create the entity. This context is added to the transaction.
		// See Txn interface
		Create(entity Entity) (opeWrap, error)

		// Close close resources
		Close() error
	}

	// Txn defines the transaction operations
	Txn interface {
		// Find checks if exists the keyValue. If it is found does DoFound or else does DoNotFound into commit
		Find(keyValue string)

		// DoFound saves the operations to transaction if it is found into commit
		DoFound(ope opeWrap)

		// DoNotFound saves the operations to transaction if it is not found into commit
		DoNotFound(ope opeWrap)

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
