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

package handler

import (
	"encoding/json"
	nethttp "net/http"
	"testing"

	"github.com/carisa/api/internal/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/api/internal/instance"

	"github.com/carisa/api/internal/http"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestInstanceHandler_Create(t *testing.T) {
	e := echo.New()
	cnt, handler, mng := newHandlerFaked(t)
	defer mng.Close()
	defer http.Close(cnt.Log, e)

	instJSON := `"name":"name","description":"desc"`
	rec, ctx := http.MockHTTP(e, "/api/instances", strings.Concat("{", instJSON, "}"))
	err := handler.Create(ctx)
	if assert.NoError(t, err) {
		assert.Contains(t, rec.Body.String(), instJSON, "Instance created")
		var inst instance.Instance
		errJ := json.NewDecoder(rec.Body).Decode(&inst)
		if assert.NoError(t, errJ) {
			assert.NotEmpty(t, inst.ID.String(), "Instance created. ID no empty")
		}
	}
}

func TestInstanceHandler_CreateWithError(t *testing.T) {
	const body = `{"name":"name","description":"desc"}`

	tests := []struct {
		name   string
		body   string
		mockS  func(*storage.ErrMockCRUD)
		mockT  func(txn *storage.ErrMockTxn)
		status int
	}{
		{
			name:   "Body wrong. Bad request",
			body:   "{df",
			status: nethttp.StatusBadRequest,
		},
		{
			name:   "Creating Entity. Error creating",
			body:   body,
			mockS:  func(s *storage.ErrMockCRUD) { s.Activate("Create") },
			status: nethttp.StatusInternalServerError,
		},
		{
			name:   "Creating Entity. Error commits transactions",
			body:   body,
			mockS:  func(s *storage.ErrMockCRUD) { s.Clear() },
			mockT:  func(s *storage.ErrMockTxn) { s.Activate("Commit") },
			status: nethttp.StatusInternalServerError,
		},
	}

	e := echo.New()
	cnt, handler, store, txn := newHandlerMocked()
	defer http.Close(cnt.Log, e)

	for _, tt := range tests {
		if tt.mockS != nil {
			tt.mockS(store)
		}
		if tt.mockT != nil {
			tt.mockT(txn)
		}
		_, ctx := http.MockHTTP(e, "/api/instances", tt.body)
		err := handler.Create(ctx)
		assert.Error(t, err, tt.name)
	}
}

func newHandlerFaked(t *testing.T) (runtime.Container, Instance, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt := mock.NewContainerFake()
	srv := instance.NewService(cnt, mng.Store())
	return cnt, NewInstanceHandl(srv, cnt), mng
}

func newHandlerMocked() (runtime.Container, Instance, *storage.ErrMockCRUD, *storage.ErrMockTxn) {
	cnt, txn := mock.NewContainerMock()
	store := &storage.ErrMockCRUD{}
	srv := instance.NewService(cnt, store)
	return cnt, NewInstanceHandl(srv, cnt), store, txn
}
