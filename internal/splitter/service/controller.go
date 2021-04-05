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

package service

import (
	"time"

	"github.com/carisa/pkg/encoding"

	"go.etcd.io/etcd/clientv3"

	"github.com/carisa/internal/splitter/runtime"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/strings"
)

const loc = "splitter.controller"

// Controller implements the functionality when the splitter service starts, stops, etc.
// Each splitter must record so node information as a timestamp indicating whether
// the splitter is alive or dead.
// If the server is dead the timestamp will not be updated and therefore the platform controller
// will know that the particular splitter is dead. The splitter controller manages
// this timestamp through ticks.
// Controller also
type Controller struct {
	cnt *runtime.Container
	// Instead of using storage.CRUD I use the etcd client directly
	// because this service is critical and I need it not to escape to the heap.
	// If one day the DB will change there are not many places where etcd is used directly.
	store *clientv3.Client

	srv  server
	tick ticks
	cons consumption

	notifyStop chan struct{}
}

// NewController builds a Controller
func NewController(cnt *runtime.Container, store *clientv3.Client) Controller {
	return Controller{
		cnt:        cnt,
		store:      store,
		tick:       newTicks(),
		cons:       newConsumption(cnt.RenewConsumptionInSecs),
		srv:        newServer(),
		notifyStop: make(chan struct{}),
	}
}

// Start starts the splitter server and registers the ticks and server information
func (c *Controller) Start() {
	splitterID := c.srv.id.String()

	c.cnt.Log.Info2(
		loc,
		"starting splitter",
		logging.String("splitter", splitterID),
		logging.String("ticks", c.tick.tstring()))

	ctx, cancel := c.cnt.StoreWithTimeout()
	_, err := c.store.KV.Put(ctx, c.keyTick(), splitterID)
	cancel()
	if err != nil {
		c.cnt.Log.Panic1(
			strings.Concat("starting splitter. error saving ticks. ", err.Error()),
			loc,
			logging.String("splitter", c.srv.id.String()))
	}

	go c.renewHeartbeat()
}

// renewHeartbeat renews the timestamp and the consumption (memory + cpu) of the splitter
// each runtime.Config.renewHeartbeatInSecs seconds
func (c *Controller) renewHeartbeat() {
	for {
		select {
		case <-c.notifyStop: // The stop service requests to terminate
			close(c.notifyStop)
			return
		case <-time.After(c.cnt.RenewHeartbeatInSecs * time.Second):
			c.updateTimestamp()
			c.updateConsumption()
		}
	}
}

// updateTimestamp updates timestamp of the splitter into db
func (c *Controller) updateTimestamp() {
	splitterID := c.srv.id.String()
	key := c.keyTick()
	c.tick.renew()
	newKey := c.keyTick()

	c.cnt.Log.Debug3(
		loc,
		"updating heartbeat of the splitter",
		logging.String("splitter", splitterID),
		logging.String("actual tick", key),
		logging.String("new tick", newKey))

	ctx, cancel := c.cnt.StoreWithTimeout()
	txn := c.store.KV.Txn(ctx).If(clientv3.Compare(clientv3.ModRevision(key), ">", 0))
	put := clientv3.OpPut(newKey, splitterID)
	txn.Then(clientv3.OpDelete(key), put)
	txn.Else(put)
	_, err := txn.Commit()
	cancel()
	if err != nil {
		c.tick.undo()
		_ = c.cnt.Log.ErrWrap1(
			err,
			loc,
			"renewHeartbeat splitter. error updating ticks",
			logging.String("ticks", key))
	}
}

// updateConsumption updates the consumption (cpu + memory) of the splitter into db
func (c *Controller) updateConsumption() {
	if err := c.cons.renew(); err != nil {
		c.cnt.Log.Panic1(
			strings.Concat("updating consumption. error getting CPU. ", err.Error()),
			loc,
			logging.String("splitter", c.srv.id.String()))
	}

	if c.cons.wake() {
		splitterID := c.srv.id.String()
		key := c.keyConsumption(c.cons.pmeasure)
		newKey := c.keyConsumption(c.cons.measure())

		c.cnt.Log.Debug3(
			loc,
			"updating consumption of the splitter",
			logging.String("splitter", splitterID),
			logging.String("actual consumption", key),
			logging.String("new consumption", newKey))

		ctx, cancel := c.cnt.StoreWithTimeout()
		txn := c.store.KV.Txn(ctx).If(clientv3.Compare(clientv3.ModRevision(key), ">", 0))
		if c.cons.cpu > 80 {
			txn.Then(clientv3.OpDelete(key))
		} else {
			put := clientv3.OpPut(newKey, c.cons.reg(c.srv.id, 1024))
			txn.Then(clientv3.OpDelete(key), put)
			txn.Else(put)
		}
		_, err := txn.Commit()
		cancel()
		if err == nil {
			c.cons.saveMeasure()
		} else {
			_ = c.cnt.Log.ErrWrap1(
				err,
				loc,
				"renewHeartbeat splitter. error updating consumption",
				logging.String("consumption", key))
		}
	}
}

// Start starts the splitter server and registers the ticks and server information.
// if it was removed return true
func (c *Controller) Stop(wait bool) bool {
	if wait {
		// it requests stop the renewHeartbeat
		c.notifyStop <- struct{}{}
	}

	key := c.keyTick()
	ctx, cancel := c.cnt.StoreWithTimeout()
	res, err := c.store.KV.Delete(ctx, key)
	cancel()
	if err != nil {
		c.cnt.Log.Panic1(
			strings.Concat("stopping splitter. error removing ticks. ", err.Error()),
			loc,
			logging.String("splitter", c.srv.id.String()))
		return false
	}
	removed := res.Deleted > 0
	if !removed {
		c.cnt.Log.Warn1(
			"stopping splitter. the ticks key is not found",
			loc,
			logging.String("ticks key", key))
	}

	return removed
}

func (c *Controller) keyTick() string {
	return strings.Concat(ticksSchema, c.tick.tstring(), c.srv.id.String())
}

func (c *Controller) keyConsumption(key uint32) string {
	keyc, _ := encoding.EncodeUI32Desc(key)
	return strings.Concat(consumptionSchema, keyc, c.srv.id.String())
}
