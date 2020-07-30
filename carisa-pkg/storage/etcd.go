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
	strs "strings"
	"testing"
	"time"

	"go.etcd.io/etcd/integration"

	"github.com/carisa/pkg/strings"

	"google.golang.org/grpc"

	"github.com/carisa/pkg/logging"

	"github.com/carisa/pkg/encoding"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// EtcdConfig defines the configuration for store framework
type EtcdConfig struct {
	// DialTimeout is the timeout for failing to establish a connection in seconds. Default value: 2 seconds.
	DialTimeout uint8 `json:"dialTimeout,omitempty"`
	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive in seconds. Default value: 10 seconds.
	DialKeepAliveTime uint8 `json:"dialKeepAliveTime,omitempty"`
	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout uint8 `json:"dialKeepAliveTimeout,omitempty"`
	// RequestTimeout in seconds. Default value: 10 seconds.
	RequestTimeout uint8 `json:"requestTimeout,omitempty"`
	// Endpoints is a list of URLs.
	Endpoints []string `json:"endpoints,omitempty"`
}

// String converts endpoint list to string
func (c *EtcdConfig) EPSString() string {
	var b strs.Builder
	b.Grow(len(c.Endpoints) * 10)
	for _, str := range c.Endpoints {
		b.WriteString(str)
		b.WriteString(",")
	}
	return b.String()
}

// etcdStore defines the CRUD operations for etcd
type etcdStore struct {
	client *clientv3.Client
}

// NewEtcd builds a store to CRUD operations from client
func NewEtcd(client *clientv3.Client) CRUD {
	return &etcdStore{client}
}

// NewEtcdConfig builds a store to CRUD operations based on etcd3 from config
func NewEtcdConfig(cnf EtcdConfig) CRUD {
	client, err := clientv3.New(config(cnf))
	if err != nil {
		panic(strings.Concat("Error creating etcd client: ", err.Error()))
	}
	return &etcdStore{client: client}
}

// Done for test
func config(cnf EtcdConfig) clientv3.Config {
	var dialTimeout time.Duration
	var dialKeepAliveTime time.Duration
	var dialKeepAliveTimeout time.Duration
	var endpoints []string

	if cnf.DialTimeout == 0 {
		dialTimeout = 2 * time.Second
	} else {
		dialTimeout = time.Duration(cnf.DialTimeout) * time.Second
	}
	if cnf.DialKeepAliveTime == 0 {
		dialKeepAliveTime = 10 * time.Second
	} else {
		dialKeepAliveTime = time.Duration(cnf.DialKeepAliveTime) * time.Second
	}
	if cnf.DialKeepAliveTimeout == 0 {
		dialKeepAliveTimeout = 2 * dialTimeout
	} else {
		dialKeepAliveTimeout = time.Duration(cnf.DialKeepAliveTimeout) * time.Second
	}
	if cnf.Endpoints == nil || len(cnf.Endpoints) == 0 {
		endpoints = append(endpoints, "localhost:2379")
	} else {
		endpoints = make([]string, len(cnf.Endpoints))
		copy(endpoints, cnf.Endpoints)
	}

	dialOptions := []grpc.DialOption{
		grpc.WithBlock(), // block until the underlying connection is up
	}
	return clientv3.Config{
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    dialKeepAliveTime,
		DialKeepAliveTimeout: dialKeepAliveTimeout,
		DialOptions:          dialOptions,
		Endpoints:            endpoints,
	}
}

// Put implements storage.interface.CRUD.Put
func (s *etcdStore) Put(entity Entity) (OpeWrap, error) {
	encode, err := encoding.Encode(entity)
	if err != nil {
		return OpeWrap{},
			errors.Wrap(
				err,
				logging.Compose("unexpected encode error putting entity into etcd store",
					logging.String("Entity", entity.ToString())))
	}

	return OpeWrap{clientv3.OpPut(entity.Key(), encode)}, err
}

// Get implements storage.interface.CRUD.Get
func (s *etcdStore) Get(ctx context.Context, key string, entity Entity) (bool, error) {
	res, err := s.client.Get(ctx, key)
	if err != nil {
		return false, errWithKey(err, key, "unexpected error getting entity from etcd store")
	}

	if res.Count > 0 {
		err = encoding.DecodeByte(res.Kvs[0].Value, entity)
		if err != nil {
			return false, errWithKey(err, key, "unexpected decode error getting entity into etcd store")
		}
		return true, nil
	}
	return false, nil
}

// Exists implements storage.interface.CRUD.Exists
func (s *etcdStore) Exists(ctx context.Context, key string) (bool, error) {
	res, err := s.client.Get(ctx, key, clientv3.WithKeysOnly())
	if err != nil {
		return false, errWithKey(err, key, "unexpected getting key into etcd store")
	}
	if res.Count > 0 {
		return string(res.Kvs[0].Key) == key, nil
	}
	return false, nil
}

func errWithKey(err error, key string, msg string) error {
	return errors.Wrap(
		err,
		logging.Compose(msg, logging.String("Key", key)))
}

// Put implements storage.interface.CRUD.Close
func (s *etcdStore) Close() error {
	return s.client.Close()
}

// etcdStore defines the operations of a transaction
type etcdTxn struct {
	client     *clientv3.Client
	opeFound   [2]OpeWrap // Could use make but this avoids escape to heap and the number of operations at most is 2
	opeNoFound [2]OpeWrap
	indexF     uint8
	indexNf    uint8
	keyValue   string
}

// Exists implements storage.interface.Txn.Exists
func (txn *etcdTxn) Find(keyValue string) {
	txn.keyValue = keyValue
}

// DoFound implements storage.interface.Txn.DoFound
func (txn *etcdTxn) DoFound(ope OpeWrap) {
	if txn.indexF > 1 {
		panic("the transaction cannot have more than 2 operations for found")
	}
	txn.opeFound[txn.indexF] = ope
	txn.indexF++
}

// DoNotFound implements storage.interface.Txn.DoNotFound
func (txn *etcdTxn) DoNotFound(ope OpeWrap) {
	if txn.indexNf > 1 {
		panic("the transaction cannot have more than 2 operations for not found")
	}
	txn.opeNoFound[txn.indexNf] = ope
	txn.indexNf++
}

// Commit implements storage.interface.Txn.Commit
func (txn *etcdTxn) Commit(ctx context.Context) (bool, error) {
	if txn.indexF == 0 && txn.indexNf == 0 {
		panic("there aren't operations")
	}
	if len(txn.keyValue) == 0 {
		panic("there isn't condition")
	}

	tx := txn.client.KV.Txn(ctx)

	if txn.indexF > 0 && txn.indexNf > 0 {
		tx = txn.ifThen(tx, ">", txn.opeFound, txn.indexF)
		if txn.indexNf == 1 {
			tx = tx.Else(txn.opeNoFound[0].opeEtcd)
		} else {
			tx = tx.Else(txn.opeNoFound[0].opeEtcd, txn.opeNoFound[1].opeEtcd)
		}
	} else {
		if txn.indexF > 0 {
			// > 0 means that the key has been found
			tx = txn.ifThen(tx, ">", txn.opeFound, txn.indexF)
		}
		if txn.indexNf > 0 {
			// = 0 means that the key has not been found
			tx = txn.ifThen(tx, "=", txn.opeNoFound, txn.indexNf)
		}
	}

	result, err := tx.Commit()

	if err != nil {
		return false, err
	}
	return result.Succeeded, nil
}

func (txn *etcdTxn) ifThen(tx clientv3.Txn, compare string, opes [2]OpeWrap, index uint8) clientv3.Txn {
	tx = tx.If(clientv3.Compare(clientv3.ModRevision(txn.keyValue), compare, 0))
	if index == 1 {
		return tx.Then(opes[0].opeEtcd)
	}
	return tx.Then(opes[0].opeEtcd, opes[1].opeEtcd)
}

type etcdIntegra struct {
	integra *integration.ClusterV3
	t       *testing.T
}

func NewEctdIntegra(t *testing.T) Integration {
	return &etcdIntegra{
		integra: integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1}),
		t:       t,
	}
}

func (e *etcdIntegra) Store() CRUD {
	return NewEtcd(e.integra.RandClient())
}

func (e *etcdIntegra) Close() {
	e.integra.Terminate(e.t)
}
