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

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestInstanceHandler_Create(t *testing.T) {
	h := mock.HTTPMock()
	cnt, handlers, _, mng := newHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	instJSON := `"name":"name","description":"desc"`
	rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/instances", strings.Concat("{", instJSON, "}"), nil)
	err := handlers.InstHandler.Create(ctx)

	if assert.NoError(t, err) {
		assert.Contains(t, rec.Body.String(), instJSON, "Created")
		assert.Equal(t, nethttp.StatusCreated, rec.Code, "Http status")
		var inst instance.Instance
		errJ := json.NewDecoder(rec.Body).Decode(&inst)
		if assert.NoError(t, errJ) {
			assert.NotEmpty(t, inst.ID.String(), "ID no empty")
		}
	}
}

func TestInstanceHandler_CreateWithError(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		mockOper func(txn *storage.ErrMockCRUDOper)
		status   int
	}{
		{
			name:   "Body wrong. Bad request",
			body:   "{df",
			status: nethttp.StatusBadRequest,
		},
		{
			name:   "Descriptor validation. Bad request",
			body:   `{"name":"","description":"desc"}`,
			status: nethttp.StatusBadRequest,
		},
		{
			name:     "Creating the instance. Error creating",
			body:     `{"name":"name","description":"desc"}`,
			mockOper: func(s *storage.ErrMockCRUDOper) { s.Activate("Create") },
			status:   nethttp.StatusInternalServerError,
		},
	}

	h := mock.HTTPMock()
	cnt, handlers, crud := newHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.mockOper != nil {
			tt.mockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/instances", tt.body, nil)
		err := handlers.InstHandler.Create(ctx)

		assert.Equal(t, tt.status, err.(*echo.HTTPError).Code, tt.name)
		assert.Error(t, err, tt.name)
	}
}

func TestInstanceHandler_Put(t *testing.T) {
	h := mock.HTTPMock()
	cnt, handlers, _, mng := newHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	params := map[string]string{"id": "12345678901234567890"}

	tests := []struct {
		instJSON string
		status   int
	}{
		{
			instJSON: `"name":"name","description":"desc"`,
			status:   nethttp.StatusCreated,
		},
		{
			instJSON: `"name":"name1","description":"desc"`,
			status:   nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/instances",
			strings.Concat("{", tt.instJSON, "}"),
			params)

		err := handlers.InstHandler.Put(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, tt.status, rec.Code, "Http status")
			assert.Contains(t, rec.Body.String(), tt.instJSON, "Put")
		}
	}
}

func TestInstanceHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": "12345678901234567890"}

	tests := []struct {
		name     string
		params   map[string]string
		body     string
		mockOper func(txn *storage.ErrMockCRUDOper)
		status   int
	}{
		{
			name:   "Body wrong. Bad request",
			params: params,
			body:   "{df",
			status: nethttp.StatusBadRequest,
		},
		{
			name:   "ID validation. Bad request",
			params: map[string]string{"i": ""},
			body:   `{"name":"name","description":"desc"}`,
			status: nethttp.StatusBadRequest,
		},
		{
			name:   "Descriptor validation. Bad request",
			params: params,
			body:   `{"name":"name","description":""}`,
			status: nethttp.StatusBadRequest,
		},
		{
			name:     "Putting the Instance. Error putting",
			params:   params,
			body:     `{"name":"name","description":"desc"}`,
			mockOper: func(s *storage.ErrMockCRUDOper) { s.Activate("Put") },
			status:   nethttp.StatusInternalServerError,
		},
	}

	h := mock.HTTPMock()
	cnt, handlers, crud := newHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.mockOper != nil {
			tt.mockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/instances", tt.body, tt.params)
		err := handlers.InstHandler.Put(ctx)

		assert.Equal(t, tt.status, err.(*echo.HTTPError).Code, tt.name)
		assert.Error(t, err, tt.name)
	}
}

func TestInstanceHandler_Get(t *testing.T) {
	h := mock.HTTPMock()
	cnt, handlers, srv, mng := newHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	inst := instance.NewInstance()
	inst.Name = "name"
	inst.Desc = "desc"

	created, err := srv.Create(&inst)
	if assert.NoError(t, err) {
		assert.True(t, created, "Instance created")

		tests := []struct {
			params map[string]string
			status int
		}{
			{
				params: map[string]string{"id": inst.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				params: map[string]string{"id": "12345678901234567890"},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/instances/:id", "", tt.params)
			err := handlers.InstHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(t, rec.Body.String(), `"name":"name","description":"desc"`, "Get instance")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestInstanceHandler_GetWithError(t *testing.T) {
	tests := []struct {
		name     string
		param    map[string]string
		mockOper func(txn *storage.ErrMockCRUDOper)
		status   int
	}{
		{
			name:   "Param not found. Bad request",
			param:  map[string]string{"i": ""},
			status: nethttp.StatusBadRequest,
		},
		{
			name:     "Get error. Internal server error",
			param:    map[string]string{"id": "12345678901234567890"},
			mockOper: func(s *storage.ErrMockCRUDOper) { s.Store().(*storage.ErrMockCRUD).Activate("Get") },
			status:   nethttp.StatusInternalServerError,
		},
	}

	h := mock.HTTPMock()
	cnt, handlers, crud := newHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.mockOper != nil {
			tt.mockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/instances/:id", "", tt.param)
		err := handlers.InstHandler.Get(ctx)

		assert.Equal(t, tt.status, err.(*echo.HTTPError).Code, tt.name)
		assert.Error(t, err, tt.name)
	}
}

func newHandlerFaked(t *testing.T) (*runtime.Container, Handlers, instance.Service, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt := mock.NewContainerFake()
	crud := storage.NewCrudOperation(mng.Store(), cnt.Log, storage.NewTxn)
	srv := instance.NewService(cnt, crud)
	hands := Handlers{InstHandler: NewInstanceHandl(srv, cnt)}
	return cnt, hands, srv, mng
}

func newHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	srv := instance.NewService(cnt, crud)
	hands := Handlers{InstHandler: NewInstanceHandl(srv, cnt)}
	return cnt, hands, crud
}
