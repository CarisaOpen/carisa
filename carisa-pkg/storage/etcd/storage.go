/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 */

package etcd

import (
	"github.com/carisa/pkg/storage"
	"github.com/coreos/etcd/clientv3"
)

// Store defines the CRUD operations for etcd
type Store struct {
	info *storage.KVMetadata
}

// Create implements storage.interface.CRUD.Create
func (s *Store) Create(entity interface{}, txn storage.Txn) {
	//txn.Do(clientv3.OpPut(s.info.GetKey(entity),  entity))
}
