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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/splitter/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestController_Start(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	defer mng.Close()

	ctrl.Start()

	_, srvID, err := mng.Store().GetRaw(context.TODO(), ctrl.keyTick())
	if assert.NoError(t, err) {
		assert.Equal(t, ctrl.srv.id.String(), srvID)
	}
}

func TestController_StartWithError(t *testing.T) {
	ctrl, txnMock, mng := newControllerMock(t)
	defer mng.Close()
	txnMock.Activate("Commit")

	assert.Panics(t, func() { ctrl.Start() })
}

func TestController_RenewHeartbeat(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	ctrl.cnt.RenewHeartbeatInSecs = 1
	defer mng.Close()

	pticks := ctrl.tick
	ctrl.Start()

	time.Sleep(2 * time.Second)
	ctrl.notifyStop <- struct{}{}
	<-ctrl.notifyStop // Wait

	assert.Equal(t, pticks.timeStamp, ctrl.tick.previousTimeStamp)
	exists, err := mng.Store().Exists(context.TODO(), strings.Concat(pticks.tstring(), ctrl.srv.id.String()))
	if assert.NoError(t, err) {
		assert.False(t, exists, "Previous tick")
		_, srvID, err := mng.Store().GetRaw(context.TODO(), ctrl.keyTick())
		if assert.NoError(t, err) {
			assert.Equal(t, ctrl.srv.id.String(), srvID, "Actual tick")
		}
	}
}

func TestController_Stop(t *testing.T) {
	ctrl, mng := newControllerFaked(t)
	ctrl.cnt.RenewHeartbeatInSecs = 1
	defer mng.Close()

	ctrl.Start()
	removed := ctrl.Stop(true)

	assert.True(t, removed)
}

func TestController_StopWithError(t *testing.T) {
	ctrl, txnMock, mng := newControllerMock(t)
	defer mng.Close()
	txnMock.Activate("Commit")

	assert.Panics(t, func() { ctrl.Stop(false) })
}

func newControllerFaked(t *testing.T) (Controller, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt := mock.NewContainerFake()
	return NewController(cnt, mng.Store()), mng
}

func newControllerMock(t *testing.T) (Controller, *storage.ErrMockTxn, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt, txnMock := mock.NewContainerMock()
	return NewController(cnt, &storage.ErrMockCRUD{}), txnMock, mng
}
