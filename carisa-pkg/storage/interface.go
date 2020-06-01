/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and  limitations under the License.
 */

package storage

// Entity defines entity context
type (
	Entity interface {
		// ToString convert entity to string
		ToString() string
	}
	// KVMetadata defines information of the entity for doing CRUD operations.
	// KVMetadata just works with key-value platforms
	KVMetadata interface {
		// GetKey gets the key
		GetKey(entity Entity) string
	}
)

type (
	// CRUD defines the CRUD operations
	CRUD interface {
		// Create creates the context to create the entity. This context is added to the transaction.
		// See Txn interface
		Create()
	}
	// Txn defines the transaction operations
	Txn interface {
		// Do save the operation to transaction
		Do(ope bag)
	}
)
