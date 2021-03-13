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
	"strconv"
	"time"

	"github.com/carisa/internal/splitter/runtime"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
)

const loc = "splitter.controller"

// Controller implements the functionality when the splitter service starts, stops, etc.
// Each splitter must record so node information as a timestamp indicating whether
// the splitter is alive or dead.
// If the server is dead the timestamp will not be updated and therefore the platform controller
// will know that the particular splitter is dead. The splitter controller manages
// this timestamp through ticks.
type Controller struct {
	cnt   *runtime.Container
	store storage.CRUD

	srv  server
	tick ticks
	cons consumption

	notifyStop chan struct{}
}

// NewController builds a Controller
func NewController(cnt *runtime.Container, data storage.CRUD) Controller {
	return Controller{
		cnt:        cnt,
		store:      data,
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

	txn := c.cnt.TxnF(c.store)

	key := c.keyTick()
	txn.Find(key)
	txn.DoNotFound(c.store.PutRaw(key, splitterID))
	ctx, cancel := c.cnt.StoreWithTimeout()
	inserted, err := txn.Commit(ctx)
	cancel()

	if err != nil {
		c.cnt.Log.Panic1(
			strings.Concat("starting splitter. error saving ticks. ", err.Error()),
			loc,
			logging.String("splitter", c.srv.id.String()))
	}
	if !inserted {
		c.cnt.Log.Panic1(
			"starting splitter. the ticks key already exists",
			loc,
			logging.String("ticks key", key))
	}

	go c.renewHeartbeat()
}

// renewHeartbeat renews the timestamp and the consumption (memory + cpu) of the splitter
// each runtime.Config.renewHeartbeatInSecs seconds
func (c *Controller) renewHeartbeat() {
	txn := c.cnt.TxnF(c.store)

	for {
		select {
		case <-c.notifyStop: // The stop service requests to terminate
			close(c.notifyStop)
			return
		case <-time.After(c.cnt.RenewHeartbeatInSecs * time.Second):
			c.updateTimestamp(txn)
			c.updateConsumption(txn)
			txn.Clear()
		}
	}
}

// updateTimestamp updates timestamp of the splitter into db
func (c *Controller) updateTimestamp(txn storage.Txn) {
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

	txn.Find(key)
	txn.DoFound(c.store.Remove(key))
	put := c.store.PutRaw(newKey, splitterID)
	txn.DoFound(put)
	txn.DoNotFound(put)
	ctx, cancel := c.cnt.StoreWithTimeout()
	_, err := txn.Commit(ctx)
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

// updateTimestamp updates the consumption (cpu + memory) of the splitter into db
func (c *Controller) updateConsumption(txn storage.Txn) {
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

		txn.Find(key)
		txn.DoFound(c.store.Remove(key))
		put := c.store.PutRaw(newKey, "1024")
		txn.DoFound(put)
		txn.DoNotFound(put)
		ctx, cancel := c.cnt.StoreWithTimeout()
		_, err := txn.Commit(ctx)
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
	txn := storage.NewTxn(c.store)
	txn.Find(key)
	txn.DoFound(c.store.Remove(key))
	ctx, cancel := c.cnt.StoreWithTimeout()
	removed, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		c.cnt.Log.Panic1(
			strings.Concat("stopping splitter. error removing ticks. ", err.Error()),
			loc,
			logging.String("splitter", c.srv.id.String()))
	}
	if !removed {
		c.cnt.Log.Warn1(
			"stopping splitter. the ticks key is not found",
			loc,
			logging.String("ticks key", key))
	}

	return removed
}

func (c *Controller) keyTick() string {
	return strings.Concat(c.tick.tstring(), c.srv.id.String())
}

func (c *Controller) keyConsumption(key int) string {
	return strings.Concat(strconv.Itoa(key), c.srv.id.String())
}
