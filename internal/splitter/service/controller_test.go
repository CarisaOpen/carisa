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
	"context"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"

	"go.etcd.io/etcd/integration"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/splitter/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestController_Start(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	defer mng.Terminate(t)

	ctrl.Start()

	res, err := ctrl.store.Get(context.TODO(), ctrl.keyTick())
	if assert.NoError(t, err) {
		assert.Equal(t, ctrl.srv.id.String(), string(res.Kvs[0].Value))
	}
}

func TestController_RenewHeartbeat(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	defer mng.Terminate(t)
	ctrl.cnt.RenewHeartbeatInSecs = 1

	pticks := ctrl.tick
	ctrl.Start()

	time.Sleep(2 * time.Second)
	ctrl.notifyStop <- struct{}{}

	assert.Equal(t, pticks.timeStamp, ctrl.tick.previousTimeStamp, "Timestamp")
	res, err := ctrl.store.Get(context.TODO(), strings.Concat(pticks.tstring(), ctrl.srv.id.String()), clientv3.WithKeysOnly())
	if assert.NoError(t, err) {
		assert.True(t, res.Count == 0, "Previous tick")
		res, err = ctrl.store.Get(context.TODO(), ctrl.keyTick())
		if assert.NoError(t, err) {
			assert.Equal(t, ctrl.srv.id.String(), string(res.Kvs[0].Value), "Actual tick")
		}
	}
}

func TestController_RenewConsumption(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	defer mng.Terminate(t)
	ctrl.cnt.RenewHeartbeatInSecs = 1

	ctrl.Start()
	time.Sleep(1 * time.Second)
	ctrl.notifyStop <- struct{}{}

	res, err := ctrl.store.Get(context.TODO(), ctrl.keyConsumption(ctrl.cons.pmeasure))
	if assert.NoError(t, err) {
		assert.NotEmpty(t, res.Kvs[0].Value)
	}
}
func TestController_Stop(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	ctrl.cnt.RenewHeartbeatInSecs = 1
	defer mng.Terminate(t)

	ctrl.Start()
	removed := ctrl.Stop(true)

	assert.True(t, removed)
}

func newControllerFaked(t *testing.T) (Controller, *integration.ClusterV3) {
	mng := storage.IntegraEtcd(t)
	cnt := mock.NewContainerFake()
	cnt.RenewConsumptionInSecs = 1
	return NewController(cnt, mng.RandClient()), mng
}
