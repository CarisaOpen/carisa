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
	"testing"
	"time"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/encoding"

	"go.etcd.io/etcd/clientv3"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/integration"
)

type EntityTest struct {
	Prop1 string
	Prop2 int8
}

func (e *EntityTest) ToString() string {
	return e.Prop1
}

func (e *EntityTest) Key() string {
	return e.Prop1
}

func TestEtcd_Config(t *testing.T) {
	tests := []struct {
		name string
		s    EtcdConfig
		t    clientv3.Config
	}{
		{
			name: "Configuring default(nil) etcd endpoints.",
			s: EtcdConfig{
				DialTimeout:          0,
				DialKeepAliveTime:    0,
				DialKeepAliveTimeout: 0,
				Endpoints:            nil,
			},
			t: clientv3.Config{
				DialTimeout:          2 * time.Second,
				DialKeepAliveTime:    10 * time.Second,
				DialKeepAliveTimeout: 2 * 2 * time.Second,
				Endpoints:            []string{"localhost:2379"},
			},
		},
		{
			name: "Configuring default([]) etcd endpoints.",
			s: EtcdConfig{
				DialTimeout:          0,
				DialKeepAliveTime:    0,
				DialKeepAliveTimeout: 0,
				Endpoints:            []string{},
			},
			t: clientv3.Config{
				DialTimeout:          2 * time.Second,
				DialKeepAliveTime:    10 * time.Second,
				DialKeepAliveTimeout: 2 * 2 * time.Second,
				Endpoints:            []string{"localhost:2379"},
			},
		},
		{
			name: "Configuring etcd",
			s: EtcdConfig{
				DialTimeout:          1,
				DialKeepAliveTime:    2,
				DialKeepAliveTimeout: 3,
				Endpoints:            []string{"server"},
			},
			t: clientv3.Config{
				DialTimeout:          1 * time.Second,
				DialKeepAliveTime:    2 * time.Second,
				DialKeepAliveTimeout: 3 * time.Second,
				Endpoints:            []string{"server"},
			},
		},
	}
	for _, tt := range tests {
		r := config(tt.s)
		assert.Equal(t, tt.t.DialTimeout, r.DialTimeout, strings.Concat(tt.name, "DialTimeout"))
		assert.Equal(t, tt.t.DialKeepAliveTime, r.DialKeepAliveTime, strings.Concat(tt.name, "DialKeepAliveTime"))
		assert.Equal(t, tt.t.DialKeepAliveTimeout, r.DialKeepAliveTimeout, strings.Concat(tt.name, "DialKeepAliveTimeout"))
		assert.Equal(t, tt.t.Endpoints, r.Endpoints, strings.Concat(tt.name, "Endpoints"))
	}
}

func TestEtcdConfig_EPSString(t *testing.T) {
	c := EtcdConfig{
		Endpoints: []string{"localhost:2378", "localhost:2379"},
	}
	assert.Equal(t, "localhost:2378,localhost:2379,", c.EPSString())
}

func TestEtcd_Create(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)

	tests := []struct {
		name string
		e    []*EntityTest
	}{
		{
			name: "Creating 1 key.",
			e: []*EntityTest{
				{
					Prop1: "key",
					Prop2: 1,
				},
			},
		},
		{
			name: "Creating 2 keys.",
			e: []*EntityTest{
				{
					Prop1: "key1",
					Prop2: 1,
				},
				{
					Prop1: "key2",
					Prop2: 1,
				},
			},
		},
		{
			name: "Creating 3 keys.",
			e: []*EntityTest{
				{
					Prop1: "key3",
					Prop2: 1,
				},
				{
					Prop1: "key4",
					Prop2: 1,
				},
				{
					Prop1: "key5",
					Prop2: 1,
				},
			},
		},
		{
			name: "Creating 4 keys.",
			e: []*EntityTest{
				{
					Prop1: "key6",
					Prop2: 1,
				},
				{
					Prop1: "key7",
					Prop2: 2,
				},
				{
					Prop1: "key8",
					Prop2: 3,
				},
				{
					Prop1: "key9",
					Prop2: 4,
				},
			},
		},
	}

	for _, tt := range tests {
		txn := NewTxn(store)
		txn.Find(tt.e[0].Prop1)

		for _, e := range tt.e {
			create, err := store.Put(e)
			if assert.NoErrorf(t, err, "%s Put failed: %v. Entity: %v", tt.name, err, e.Prop1) {
				txn.DoNotFound(create)
			}
		}

		ok, errC := txn.Commit(ctx)
		if assert.NoErrorf(t, errC, "%s Commit failed: %v", tt.name, errC) {
			assert.True(t, ok, strings.Concat(tt.name, "Entity found"))
			for _, e := range tt.e {
				var er EntityTest
				_, errG := store.Get(ctx, e.Prop1, &er)
				if assert.NoErrorf(t, errC, "%s Get failed: %v. Entity: $s", tt.name, errG, e.Prop1) {
					assert.Equalf(t, e, &er, "%v Entity '%v' not created", tt.name, e.Prop1)
				}
			}
		}
	}
}

func TestEtcd_Update(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	tests := []struct {
		name string
		e    []*EntityTest
	}{
		{
			name: "Updating.",
			e: []*EntityTest{
				{
					Prop1: "key",
					Prop2: 1,
				},
			},
		},
		{
			name: "Updating many keys.",
			e: []*EntityTest{
				{
					Prop1: "key1",
					Prop2: 1,
				},
				{
					Prop1: "key2",
					Prop2: 2,
				},
				{
					Prop1: "key3",
					Prop2: 3,
				},
				{
					Prop1: "key4",
					Prop2: 4,
				},
			},
		},
	}

	for _, tt := range tests {
		txn := NewTxn(store)
		txn.Find(tt.e[0].Prop1)

		for _, e := range tt.e {
			_, errp := client.KV.Put(ctx, e.Prop1, "1")
			if assert.NoErrorf(t, errp, "%s Put KV failed: %v. Entity: %s", tt.name, errp, e.Prop1) {
				update, err := store.Put(e)
				if assert.NoErrorf(t, err, "%s Put failed: %v. Entity: %s", tt.name, err, e.Prop1) {
					txn.DoFound(update)
				}
			}
		}

		ok, errC := txn.Commit(ctx)
		if assert.NoErrorf(t, errC, "%s Commit failed: %v", tt.name, errC) {
			assert.True(t, ok, strings.Concat(tt.name, "Entity not found"))
			for _, e := range tt.e {
				var er EntityTest
				_, errG := store.Get(ctx, e.Prop1, &er)
				if assert.NoErrorf(t, errC, "%s Get failed: %v. Entity: $s", tt.name, errG, e.Prop1) {
					assert.Equalf(t, e, &er, e.Prop1, "%s Entity '%s' not created", tt.name, e.Prop1)
				}
			}
		}
	}
}

func TestEtcd_Put(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)

	tests := []struct {
		name  string
		e     *EntityTest
		found bool
	}{
		{
			name: "Creating.",
			e: &EntityTest{
				Prop1: "key",
				Prop2: 1,
			},
			found: false,
		},
		{
			name: "Updating.",
			e: &EntityTest{
				Prop1: "key",
				Prop2: 2,
			},
			found: true,
		},
	}

	txn := NewTxn(store)
	txn.Find(tests[0].e.Prop1)

	putnf, err := store.Put(tests[0].e)
	if assert.NoErrorf(t, err, "%Put to create failed: %v. Entity: %s", err, tests[0].e.Prop1) {
		putf, err := store.Put(tests[1].e)
		if assert.NoErrorf(t, err, "Put to Update failed: %v. Entity: %s", err, tests[1].e.Prop1) {
			txn.DoNotFound(putnf) // Create
			txn.DoFound(putf)     // Update

			for _, tt := range tests {
				found, err := txn.Commit(ctx)
				if assert.NoErrorf(t, err, "%s Commit failed: %v", tt.name, err) {
					assert.Equal(t, found, tt.found, strings.Concat(tt.name, "Commit result"))
					var er EntityTest
					_, err := store.Get(ctx, tt.e.Prop1, &er)
					if assert.NoErrorf(t, err, "%s Get failed: %v. Entity: $s", tt.name, err, tt.e.Prop1) {
						assert.Equalf(t, tt.e, &er, tt.e.Prop1, "%s Entity '%s' not created", tt.name, tt.e.Prop1)
					}
				}
			}
		}
	}
}

func TestEtcd_PutRaw(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	key := "key"
	value := "value"

	txn := NewTxn(store)
	txn.Find(key)
	txn.DoNotFound(store.PutRaw(key, value))
	_, err := txn.Commit(ctx)
	if assert.NoError(t, err) {
		res, err := client.Get(ctx, key)
		if assert.NoError(t, err) {
			assert.Equal(t, key, string(res.Kvs[0].Key), "Key")
			assert.Equal(t, value, string(res.Kvs[0].Value), "Value")
		}
	}
}

func TestEtcd_Get(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	e := &EntityTest{
		Prop1: "key",
		Prop2: 1,
	}

	result := []struct {
		name  string
		key   string
		found bool
	}{
		{
			name:  "Getting. Found.",
			key:   e.Prop1,
			found: true,
		},
		{
			name:  "Getting. Not found.",
			key:   "key2",
			found: false,
		},
	}

	entitye, err := encoding.Encode(e)
	if assert.NoErrorf(t, err, "Coding entity") {
		_, err := client.Put(ctx, e.Prop1, entitye)
		if assert.NoErrorf(t, err, "Put entity") {
			for _, tt := range result {
				var entityg EntityTest
				found, err := store.Get(ctx, tt.key, &entityg)
				if assert.NoErrorf(t, err, strings.Concat(tt.name, "Get entity")) {
					assert.Equal(t, tt.found, found, strings.Concat(tt.name, "Entity found"))
					if found {
						assert.Equal(t, e, &entityg, strings.Concat(tt.name, "Get result"))
					}
				}
			}
		}
	}
}

func TestEtcd_Remove(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	result := []struct {
		name    string
		removed bool
	}{
		{
			name:    "Removing. Ok.",
			removed: true,
		},
		{
			name:    "Removing. Not found.",
			removed: false,
		},
	}

	const key = "key"

	txn := NewTxn(store)
	txn.Find(key)

	_, err := client.Put(ctx, key, "value")
	if assert.NoError(t, err) {
		for _, tt := range result {
			txn.DoFound(store.Remove(key))
			ok, err := txn.Commit(ctx)
			if assert.NoErrorf(t, err, strings.Concat(tt.name, "Remove entity")) {
				assert.Equal(t, tt.removed, ok)
				exists, err := store.Exists(ctx, key)
				if assert.NoErrorf(t, err, strings.Concat(tt.name, "Entity exists")) {
					assert.False(t, exists)
				}
			}
		}
	}
}

func TestEtcd_Exists(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	result := []struct {
		name  string
		key   string
		found bool
	}{
		{
			name:  "Found.",
			key:   "key",
			found: true,
		},
		{
			name:  "Not found.",
			key:   "key1",
			found: false,
		},
	}

	_, err := client.Put(ctx, "key", "value")
	if assert.NoErrorf(t, err, "Put entity") {
		for _, tt := range result {
			found, err := store.Exists(ctx, tt.key)
			if assert.NoErrorf(t, err, strings.Concat(tt.name, "Exists entity")) {
				assert.Equal(t, tt.found, found, strings.Concat(tt.name, "Entity found"))
			}
		}
	}
}

func TestEtcd_StartKey(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	samples := []EntityTest{
		{
			Prop1: "key11",
			Prop2: 11,
		},
		{
			Prop1: "key1",
			Prop2: 1,
		},
		{
			Prop1: "key0",
			Prop2: 0,
		},
		{
			Prop1: "ky1",
			Prop2: 1,
		},
	}

	if !sampling(ctx, t, samples, client) {
		return
	}

	tests := []struct {
		name string
		key  string
		top  int
		res  []EntityTest
	}{
		{
			name: "Top 0.",
			key:  "ky",
			top:  0,
			res: []EntityTest{
				{
					Prop1: "ky1",
					Prop2: 1,
				},
			},
		},
		{
			name: "Not found.",
			key:  "y1",
			top:  5,
			res:  []EntityTest{},
		},
		{
			name: "Found. < 10. Top 10.",
			key:  "key1",
			top:  10,
			res: []EntityTest{
				{
					Prop1: "key1",
					Prop2: 1,
				},
				{
					Prop1: "key11",
					Prop2: 11,
				},
			},
		},
		{
			name: "Found. > 10. Top 1.",
			key:  "key",
			top:  1,
			res: []EntityTest{
				{
					Prop1: "key0",
					Prop2: 0,
				},
			},
		},
		{
			name: "Other key.",
			key:  "ky1",
			top:  2,
			res: []EntityTest{
				{
					Prop1: "ky1",
					Prop2: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		res, err := store.StartKey(ctx, tt.key, tt.top, func() Entity { return &EntityTest{} })
		if assert.NoErrorf(t, err, "StartKey") {
			assert.Equal(t, len(tt.res), len(res), strings.Concat(tt.name, "Count"))
			for i, r := range res {
				assert.Equal(t, &tt.res[i], r, strings.Concat(tt.name, "StartKey result"))
			}
		}
	}
}

func TestEtcd_Range(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	samples := []EntityTest{
		{
			Prop1: "key1",
			Prop2: 1,
		},
		{
			Prop1: "key0",
			Prop2: 0,
		},
		{
			Prop1: "key2",
			Prop2: 2,
		},
		{
			Prop1: "ky1",
			Prop2: 1,
		},
	}

	if !sampling(ctx, t, samples, client) {
		return
	}

	tests := []struct {
		name string
		skey string
		ekey string
		top  int
		res  []EntityTest
	}{
		{
			name: "Not found.",
			skey: "y1",
			ekey: "y",
			top:  5,
			res:  []EntityTest{},
		},
		{
			name: "Found. < 10. Top 10.",
			skey: "key1",
			ekey: "key",
			top:  10,
			res: []EntityTest{
				{
					Prop1: "key1",
					Prop2: 1,
				},
				{
					Prop1: "key2",
					Prop2: 2,
				},
			},
		},
		{
			name: "Found. < 10. Top 3.",
			skey: "key",
			ekey: "key",
			top:  3,
			res: []EntityTest{
				{
					Prop1: "key0",
					Prop2: 0,
				},
				{
					Prop1: "key1",
					Prop2: 1,
				},
				{
					Prop1: "key2",
					Prop2: 2,
				},
			},
		},
		{
			name: "Other key.",
			skey: "key1",
			ekey: "key",
			top:  1,
			res: []EntityTest{
				{
					Prop1: "key1",
					Prop2: 1,
				},
			},
		},
		{
			name: "Top 0.",
			skey: "key2",
			ekey: "key",
			top:  0,
			res: []EntityTest{
				{
					Prop1: "key2",
					Prop2: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		res, err := store.Range(ctx, tt.skey, tt.ekey, tt.top, func() Entity { return &EntityTest{} })
		if assert.NoErrorf(t, err, "Range") {
			assert.Equal(t, len(tt.res), len(res), strings.Concat(tt.name, "Count"))
			for i, r := range res {
				assert.Equal(t, &tt.res[i], r, strings.Concat(tt.name, "Range result"))
			}
		}
	}
}

func TestEtcd_RangeRaw(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	samples := map[string]string{
		"key1": "1",
		"key0": "0",
		"key2": "2",
		"ky1":  "1",
	}

	for k, v := range samples {
		_, err := client.Put(ctx, k, v)
		if err != nil {
			assert.NoError(t, err, "Putting samples")
			return
		}
	}

	tests := []struct {
		name string
		skey string
		ekey string
		top  int
		res  map[string]string
	}{
		{
			name: "Not found.",
			skey: "y1",
			ekey: "y",
			top:  5,
			res:  map[string]string{},
		},
		{
			name: "Found. < 10. Top 10.",
			skey: "key1",
			ekey: "key",
			top:  10,
			res: map[string]string{
				"key1": "1",
				"key2": "2",
			},
		},
		{
			name: "Found. < 10. Top 3.",
			skey: "key",
			ekey: "key",
			top:  3,
			res: map[string]string{
				"key0": "0",
				"key1": "1",
				"key2": "2",
			},
		},
		{
			name: "Other key.",
			skey: "key1",
			ekey: "key",
			top:  1,
			res: map[string]string{
				"key1": "1",
			},
		},
	}

	for _, tt := range tests {
		res, err := store.RangeRaw(ctx, tt.skey, tt.ekey, tt.top)
		if assert.NoErrorf(t, err, "RangeRaw") {
			assert.Equal(t, len(tt.res), len(res), strings.Concat(tt.name, "Count"))
			assert.Equal(t, tt.res, res, strings.Concat(tt.name, "Range result"))
		}
	}
}

func sampling(ctx context.Context, t *testing.T, samples []EntityTest, client *clientv3.Client) bool {
	for _, s := range samples {
		entity, err := encoding.Encode(s)
		if err != nil {
			assert.NoError(t, err, "Coding entity for sample")
			return false
		}
		_, err = client.Put(ctx, s.Prop1, entity)
		if err != nil {
			assert.NoError(t, err)
			return false
		}
	}
	return true
}

func TestEtcd_IntegraStore(t *testing.T) {
	i := NewEctdIntegra(t)
	defer i.Close()
	assert.NotNil(t, i.Store())
}

func newStore(t *testing.T) (*integration.ClusterV3, context.Context, CRUD) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	store := NewEtcd(cluster.RandClient())
	ctx := context.TODO()
	return cluster, ctx, store
}
