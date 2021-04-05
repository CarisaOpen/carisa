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
	"bytes"

	"github.com/carisa/internal/controller/runtime"
	"github.com/carisa/pkg/logging"
	"go.etcd.io/etcd/clientv3"
)

const loc = "controller.selector"

const (
	consumptionSchema = "SC"
	// top is the numbers of splitter consumption that can read in each page
	top = 20
	// defaultEntesAvg is the default average of entes that the controller typically router
	defaultEntesAvg = 10
)

// selector prioritises the splitters
// that will be chosen for the management of each Ente
// according to the consumption of each one.
type selector struct {
	cnt *runtime.Container
	// Instead of using storage.CRUD I use the etcd client directly
	// because this service is critical and I need it not to escape to the heap.
	// If one day the DB will change there are not many places where etcd is used directly.
	store *clientv3.Client

	// entesAvg is the average of entes that the controller typically router
	entesAvg uint32
}

// newSelector builds a selector
func newSelector(cnt *runtime.Container, store *clientv3.Client) selector {
	return selector{
		cnt:      cnt,
		store:    store,
		entesAvg: defaultEntesAvg,
	}
}

// sort sorts and ranks all splitter according to CPU and memory.
// Selects the splitter with the highest memory and CPU usage.
// The data is extracted from the database and paged from top to top.
// The goal is to make maximum use of the splitter without saturating it.
// sort returns a list of splitter ranked by priority for handling entes requests
// and also carries the maximum number of entes that a splitter could contain.
func (s *selector) sort(top uint8 /*For testing*/) []splitter {
	var splitters []splitter
	var entesCounter uint32 = 0
	lastKey := []byte(consumptionSchema)
	ppag := clientv3.WithPrefix()

	// it reading the consumption order by ascending, although as etcd does not allow paging
	// in descending order, the key is coded so that applying an ascending order will return
	// the records in descending order.
	for {
		ctx, cancel := s.cnt.StoreWithTimeout()
		consumptions, err := s.store.Get(
			ctx, string(lastKey),
			clientv3.WithLimit(int64(top)),
			ppag,
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
		cancel()
		if err != nil {
			_ = s.cnt.Log.ErrWrap(
				err,
				loc,
				"error reading consumption of the splitter")
			return nil
		}

		for _, consumption := range consumptions.Kvs {
			if bytes.Equal(lastKey, consumption.Key) {
				continue
			}
			lastKey = consumption.Key
			if entesCounter > s.entesAvg {
				return splitters
			}
			consm, err := unmarshallConsumption(consumption.Value)
			if err != nil {
				_ = s.cnt.Log.ErrWrap1(
					err,
					loc,
					"error unmarshalling consumption of the splitter",
					logging.String("key", consumption.String()))
				continue
			}

			// Number of entes the splitter can handle
			if consm.mem >= s.cnt.SplitterMaxMemoryInMB {
				continue
			}
			entes := (s.cnt.SplitterMaxMemoryInMB - consm.mem) / consm.entesMemAvg

			entesCounter += entes
			splitters = append(splitters, splitter{
				id:    consm.splitterID,
				entes: entes,
			})
		}

		if len(consumptions.Kvs) != int(top) { // Eof
			break
		}
		ppag = clientv3.WithFromKey()
	}
	return splitters
}
