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

import (
	"fmt"
	"github.com/carisa/pkg/encoding"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// StoreEtcd defines the CRUD operations for etcd
type StoreEtcd struct {
	info KVMetadata
}

// Create implements storage.interface.CRUD.Create
func (s StoreEtcd) Create(entity Entity) (error, bag) {

	encode, err := encoding.Encode(entity)
	if err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf(
				"Unexpected encode error creating entity into etcd store. Entity: (%v)",
				entity.ToString())), nil
	}

	return err, bagEtcd(clientv3.OpPut(s.info.GetKey(entity), encode))
}

// StoreEtcd defines the operations of a transaction
type TxnEtcd struct {
	client *clientv3.Client
	opes   []bag
}

func NewTxnEtcd(client *clientv3.Client) TxnEtcd {
	return TxnEtcd{
		client: client,
		opes:   make([]bag, 2),
	}
}

// Do implements storage.interface.Txn.Do
func (txn *TxnEtcd) Do(ope bag, err error) error {
	txn.opes = append(txn.opes, ope)
	return err
}
