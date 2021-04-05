/*
 *  Copyright 2019-2022 the original author or authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package controller

import (
	"context"
	"testing"

	"github.com/carisa/pkg/runtime"

	"github.com/carisa/internal/controller/mock"
	"github.com/carisa/pkg/encoding"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/integration"
)

func TestSelector_sort(t *testing.T) {
	s, mng := newSelectorFaked(t)
	defer mng.Terminate(t)

	tests := []struct {
		name      string
		top       uint8
		entesAvg  uint32
		splitters []consumption
		sorted    []splitter
	}{
		{
			name:     "No paging, selecting all and mem > SplitterMaxMemoryInMB",
			top:      top,
			entesAvg: 1000,
			splitters: []consumption{
				{
					splitterID:  xid.NilID(),
					cpu:         30,
					mem:         512,
					entesMemAvg: 50,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         50,
					mem:         256,
					entesMemAvg: 20,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         60,
					mem:         2048,
					entesMemAvg: 20,
				},
			},
			sorted: []splitter{
				{
					id:    xid.NilID(),
					entes: 38,
				},
				{
					id:    xid.NilID(),
					entes: 10,
				},
			},
		},
		{
			name:     "No paging, entesAvg < total of consumption",
			top:      top,
			entesAvg: 45,
			splitters: []consumption{
				{
					splitterID:  xid.NilID(),
					cpu:         30,
					mem:         512,
					entesMemAvg: 50,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         50,
					mem:         256,
					entesMemAvg: 20,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         40,
					mem:         256,
					entesMemAvg: 5,
				},
			},
			sorted: []splitter{
				{
					id:    xid.NilID(),
					entes: 38,
				},
				{
					id:    xid.NilID(),
					entes: 10,
				},
			},
		},
		{
			name:     "Paging",
			top:      2,
			entesAvg: 1000,
			splitters: []consumption{
				{
					splitterID:  xid.NilID(),
					cpu:         30,
					mem:         256,
					entesMemAvg: 50,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         50,
					mem:         256,
					entesMemAvg: 20,
				},
				{
					splitterID:  xid.NilID(),
					cpu:         70,
					mem:         256,
					entesMemAvg: 5,
				},
			},
			sorted: []splitter{
				{
					id:    xid.NilID(),
					entes: 153,
				},
				{
					id:    xid.NilID(),
					entes: 38,
				},
				{
					id:    xid.NilID(),
					entes: 15,
				},
			},
		},
	}

	client := mng.RandClient()
	for _, tt := range tests {
		if err := samples(client, tt.splitters); err != nil {
			return
		}
		s.entesAvg = tt.entesAvg
		res := s.sort(tt.top)
		assert.Equalf(t, tt.sorted, res, tt.name)
		if err := unsamples(client); err != nil {
			return
		}
	}
}

func samples(client *clientv3.Client, splitters []consumption) error {
	for _, s := range splitters {
		sc := encoding.SimpleEncoder{}
		sc.WriteBytes(s.splitterID.Bytes())
		sc.WriteUint8(s.cpu)
		sc.WriteUint32(s.mem)
		sc.WriteUint32(s.entesMemAvg)

		if _, err := client.Put(context.TODO(), s.key(), sc.String()); err != nil {
			return err
		}
	}
	return nil
}

func unsamples(client *clientv3.Client) (err error) {
	_, err = client.Delete(context.TODO(), consumptionSchema, clientv3.WithPrefix())
	return
}

// measure can change. Look at splitter service
func (c *consumption) key() string {
	codecm, _ := encoding.EncodeUI32Desc(runtime.Meassure(c.cpu, c.mem))
	return strings.Concat(
		consumptionSchema,
		codecm,
		c.splitterID.String())
}

func newSelectorFaked(t *testing.T) (selector, *integration.ClusterV3) {
	mng := storage.IntegraEtcd(t)
	cnt := mock.NewContainerFake()
	return newSelector(cnt, mng.RandClient()), mng
}
