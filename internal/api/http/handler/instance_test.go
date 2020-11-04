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
	"fmt"
	nethttp "net/http"
	"testing"

	"github.com/carisa/internal/api/service"
	spacesmpl "github.com/carisa/internal/api/space/samples"

	"github.com/rs/xid"

	"github.com/carisa/internal/api/samples"

	"github.com/carisa/internal/api/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/api/mock"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/instance"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestInstanceHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newInstHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	instJSON := `"name":"name","description":"desc"`
	rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/instances", strings.Concat("{", instJSON, "}"), nil, nil)
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
	tests := samples.TestCreateWithError("Create")

	h := mock.HTTP()
	cnt, handlers, crud := newInstHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/instances", tt.Body, nil, nil)
		err := handlers.InstHandler.Create(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestInstanceHandler_Put(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newInstHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating instance.",
			body:   `"name":"name","description":"desc"`,
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating instance.",
			body:   `"name":"name1","description":"desc"`,
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/instances",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)

		err := handlers.InstHandler.Put(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Put"))
		}
	}
}

func TestInstanceHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := samples.TestPutWithError("Put", params)

	h := mock.HTTP()
	cnt, handlers, crud := newInstHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/instances", tt.Body, tt.Params, nil)
		err := handlers.InstHandler.Put(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestInstanceHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newInstHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	inst := instance.New()
	inst.Name = "name"
	inst.Desc = "desc"

	created, err := srv.Create(&inst)
	if assert.NoError(t, err) {
		assert.True(t, created, "Instance created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Find instance. Instance found.",
				params: map[string]string{"id": inst.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Find instance. Instance not found.",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/instances/:id", "", tt.params, nil)
			err := handlers.InstHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						`"name":"name","description":"desc"`,
						strings.Concat(tt.name, "Get instance"))
				}
				assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			}
		}
	}
}

func TestInstanceHandler_GetWithError(t *testing.T) {
	tests := samples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newInstHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/instances/:id", "", tt.Param, nil)
		err := handlers.InstHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestInstanceHandler_ListSpaces(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newInstHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, space, err := spacesmpl.CreateLink(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/instances/:id/spaces",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "name"})

		err := handlers.InstHandler.ListSpaces(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"name","spaceId":"%s"}]`, space.ID.String()),
				"List space")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestInstanceHandler_GetListSpacesError(t *testing.T) {
	tests := samples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newInstHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/instances/:id/spaces", "", tt.Param, tt.QParam)
		err := handlers.InstHandler.ListSpaces(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newInstHandlerFaked(t *testing.T) (*runtime.Container, Handlers, instance.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	srv := instance.NewService(cnt, ext, crud)
	hands := Handlers{InstHandler: NewInstanceHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newInstHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	srv := instance.NewService(cnt, ext, crud)
	hands := Handlers{InstHandler: NewInstanceHandle(srv, cnt)}
	return cnt, hands, crud
}
