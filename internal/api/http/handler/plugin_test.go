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

	psamples "github.com/carisa/internal/api/plugin/samples"

	tsamples "github.com/carisa/internal/api/samples"

	"github.com/carisa/pkg/strings"
	"github.com/labstack/echo/v4"

	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/plugin"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestPluginHandler_Create(t *testing.T) {
	cnt, handlers, _, mng := newPluginHandlerFaked(t)
	defer mng.Close()
	h := mock.HTTP()
	defer h.Close(cnt.Log)

	plugins := plugins()

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating plugin.",
			body:   `"name":"name","description":"desc"`,
			status: nethttp.StatusCreated,
		},
	}

	for _, pc := range plugins {
		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/plugins", strings.Concat("{", tt.body, "}"), nil, nil)
			err := handlers.PluginHandler.Create(ctx, pc)

			if err != nil && tt.status == err.(*echo.HTTPError).Code {
				continue
			}
			if err != nil {
				assert.Error(t, err, strings.Concat(tt.name, "Error create"))
				continue
			}
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var proto plugin.Prototype
				errj := json.NewDecoder(rec.Body).Decode(&proto)
				if assert.NoError(t, errj) {
					assert.NotEmpty(t, proto.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestPluginHandler_CreateWithError(t *testing.T) {
	cnt, handlers, crud := newPluginHandlerMocked()

	plugins := plugins()

	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	defer h.Close(cnt.Log)

	for _, pc := range plugins {
		for _, tt := range tests {
			if tt.MockOper != nil {
				tt.MockOper(crud)
			}
			_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/plugins", tt.Body, nil, nil)
			err := handlers.PluginHandler.Create(ctx, pc)

			assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
			assert.Error(t, err, tt.Name)
		}
	}
}

func TestPluginHandler_Put(t *testing.T) {
	cnt, handlers, _, mng := newPluginHandlerFaked(t)
	defer mng.Close()
	h := mock.HTTP()
	defer h.Close(cnt.Log)

	plugins := plugins()

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating plugin.",
			body:   `"name":"name","description":"desc"`,
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating plugin.",
			body:   `"name":"name1","description":"desc"`,
			status: nethttp.StatusOK,
		},
	}

	for _, pc := range plugins {
		for _, tt := range tests {
			rec, ctx := h.NewHTTP(
				nethttp.MethodPut,
				"/api/plugins",
				strings.Concat("{", tt.body, "}"),
				params,
				nil)

			err := handlers.PluginHandler.Put(ctx, pc)
			if err != nil && tt.status == err.(*echo.HTTPError).Code {
				continue
			}
			if assert.NoError(t, err, tt.name) {
				assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
				if tt.status != nethttp.StatusNotFound {
					assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Put"))
				}
			}
		}
	}
}

func TestPluginHandler_PutWithError(t *testing.T) {
	cnt, handlers, crud := newPluginHandlerMocked()

	plugins := plugins()

	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	defer h.Close(cnt.Log)

	for _, pc := range plugins {
		for _, tt := range tests {
			if tt.MockOper != nil {
				tt.MockOper(crud)
			}
			_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/plugins", tt.Body, tt.Params, nil)
			err := handlers.PluginHandler.Put(ctx, pc)

			assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
			assert.Error(t, err, tt.Name)
		}
	}
}

func TestPluginHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newPluginHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	proto := plugin.New()
	proto.Name = "pname"
	proto.Desc = "pdesc"
	created, err := srv.Create(&proto)

	if assert.NoError(t, err) {
		assert.True(t, created, "Plugin created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding proto. Ok",
				params: map[string]string{"id": proto.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding proto. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/plugins/:id", "", tt.params, nil)
			err := handlers.PluginHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						`"name":"pname","description":"pdesc"`,
						"Get proto")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestPluginHandler_GetWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newPluginHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/plugins/:id", "", tt.Param, nil)
		err := handlers.PluginHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestPluginHandler_ListPlugins(t *testing.T) {
	cnt, handlers, _, mng := newPluginHandlerFaked(t)
	defer mng.Close()
	h := mock.HTTP()
	defer h.Close(cnt.Log)

	plugins := plugins()

	for _, pc := range plugins {
		_, proto, err := psamples.CreateLinkPlugin(mng, pc)

		if assert.NoError(t, err) {
			rec, ctx := h.NewHTTP(
				nethttp.MethodGet,
				"/api/platforms/plugins",
				"",
				map[string]string{"id": xid.NilID().String()},
				map[string]string{"sname": "nameproto"})

			err := handlers.PluginHandler.ListPlugins(ctx, pc)
			if assert.NoError(t, err) {
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`[{"name":"nameproto","protoId":"%s","category":"%s"}]`, proto.Key(), string(pc)),
					"List the queries plugin")
				assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
			}
		}
	}
}

func TestPluginHandler_ListPluginsError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newPluginHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/platform/plugins", "", tt.Param, tt.QParam)
		err := handlers.PluginHandler.ListPlugins(ctx, plugin.Query)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newPluginHandlerFaked(t *testing.T) (*runtime.Container, Handlers, plugin.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	srv := plugin.NewService(cnt, ext, crud)
	hands := Handlers{PluginHandler: NewPluginHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newPluginHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	srv := plugin.NewService(cnt, ext, crud)
	hands := Handlers{PluginHandler: NewPluginHandle(srv, cnt)}
	return cnt, hands, crud
}

func plugins() []plugin.Category {
	return []plugin.Category{plugin.Query}
}
