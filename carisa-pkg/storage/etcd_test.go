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

func (e *EntityTest) GetKey() string {
	return e.Prop1
}

func TestEtcdCreate(t *testing.T) {
	cluster, ctx, store := newStore(t)
	defer cluster.Terminate(t)
	client := cluster.RandClient()

	tests := []struct {
		e []*EntityTest
	}{
		{
			e: []*EntityTest{
				{
					Prop1: "key",
					Prop2: 1,
				},
			},
		},
		{
			e: []*EntityTest{
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
	}

	for _, tt := range tests {
		txn := NewTxn(store)
		txn.Find(tt.e[0].GetKey())

		for _, e := range tt.e {
			create, err := store.Create(e)
			assert.NoErrorf(t, err, "Create failed: %v. Entity: %s", err, e.GetKey())
			txn.DoNotFound(create)
		}

		ok, errC := txn.Commit(ctx)
		assert.NoErrorf(t, errC, "Commit failed: %v", errC)
		assert.True(t, ok, "Entity found")

		for _, e := range tt.e {
			r, errG := client.KV.Get(ctx, e.GetKey())
			assert.NoErrorf(t, errC, "Get failed: %v. Entity: $s", errG, e.GetKey())
			assert.Equalf(t, string(r.Kvs[0].Key), e.GetKey(), "Entity '%s' not saved", e.GetKey())
		}
	}
}

func newStore(t *testing.T) (*integration.ClusterV3, context.Context, CRUD) {
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	store := NewEtcdStore(cluster.RandClient())
	ctx := context.TODO()
	return cluster, ctx, store
}
