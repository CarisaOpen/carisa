/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package storage

import (
	"context"
	"github.com/carisa/pkg/encoding"
	"github.com/carisa/pkg/strings"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// etcdStore defines the CRUD operations for etcd
type etcdStore struct {
	client *clientv3.Client
}

// NewEtcdStore builds a store to CRUD operations based on etcd3
func NewEtcdStore(client *clientv3.Client) CRUD {
	return &etcdStore{client}
}

// Create implements storage.interface.CRUD.Create
func (s *etcdStore) Create(entity Entity) (opeWrap, error) {

	encode, err := encoding.Encode(entity)
	if err != nil {
		return opeWrap{},
			errors.Wrap(
				err,
				strings.Concat(
					"Unexpected encode error creating entity into etcd store. Entity: ",
					entity.ToString()))
	}

	return opeWrap{clientv3.OpPut(entity.GetKey(), encode)}, err
}

// etcdStore defines the operations of a transaction
type etcdTxn struct {
	client     *clientv3.Client
	opeFound   [2]opeWrap // Could use make but this avoids escape to heap and the number of operations at most is 2
	opeNoFound [2]opeWrap
	indexF     uint8
	indexNf    uint8
	keyValue   string
}

// Exists implements storage.interface.Txn.Exists
func (txn *etcdTxn) Find(keyValue string) {
	txn.keyValue = keyValue
}

// DoFound implements storage.interface.Txn.DoFound
func (txn *etcdTxn) DoFound(ope opeWrap) {
	if txn.indexF > 1 {
		panic("The transaction cannot have more than 2 operations for found")
	}
	txn.opeFound[txn.indexF] = ope
	txn.indexF++
}

// DoNotFound implements storage.interface.Txn.DoNotFound
func (txn *etcdTxn) DoNotFound(ope opeWrap) {
	if txn.indexNf > 1 {
		panic("The transaction cannot have more than 2 operations for not found")
	}
	txn.opeNoFound[txn.indexNf] = ope
	txn.indexNf++
}

// Commit implements storage.interface.Txn.Commit
func (txn *etcdTxn) Commit(ctx context.Context) (bool, error) {

	if txn.indexF == 0 && txn.indexNf == 0 {
		panic("There aren't operations")
	}
	if len(txn.keyValue) == 0 {
		panic("There isn't condition")
	}

	tx := txn.client.KV.Txn(ctx)

	if txn.indexF > 0 && txn.indexNf > 0 {
		txn.ifThen(tx, ">", txn.opeFound, txn.indexF)
		if txn.indexNf == 1 {
			tx = tx.Else(txn.opeNoFound[0].opeEtcd)
		} else {
			tx = tx.Else(txn.opeNoFound[0].opeEtcd, txn.opeNoFound[1].opeEtcd)
		}
	} else {
		if txn.indexF > 0 {
			// > 0 means that the key has been found
			txn.ifThen(tx, ">", txn.opeFound, txn.indexF)
		}
		if txn.indexNf > 0 {
			// = 0 means that the key has not been found
			txn.ifThen(tx, "=", txn.opeNoFound, txn.indexNf)
		}
	}

	result, err := tx.Commit()

	if err != nil {
		return false, err
	}
	return result.Succeeded, nil
}

func (txn *etcdTxn) ifThen(tx clientv3.Txn, compare string, opes [2]opeWrap, index uint8) {
	tx = tx.If(clientv3.Compare(clientv3.ModRevision(txn.keyValue), compare, 0))
	if index == 1 {
		tx = tx.Then(opes[0].opeEtcd)
	} else {
		tx = tx.Then(opes[0].opeEtcd, opes[1].opeEtcd)
	}
}
