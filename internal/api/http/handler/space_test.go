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

	catsmpl "github.com/carisa/internal/api/category/samples"
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/internal/api/space"

	entesmpl "github.com/carisa/internal/api/ente/samples"

	"github.com/rs/xid"

	"github.com/carisa/internal/api/instance/samples"
	tsamples "github.com/carisa/internal/api/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestSpaceHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newSpcHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	inst, err := samples.CreateInstance(mng)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating space.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Creating space. Instance not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/spaces", strings.Concat("{", tt.body, "}"), nil, nil)
		err := handlers.SpaceHandler.Create(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var spc space.Space
				errJ := json.NewDecoder(rec.Body).Decode(&spc)
				if assert.NoError(t, errJ) {
					assert.NotEmpty(t, spc.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestSpaceHandler_CreateWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/spaces", tt.Body, nil, nil)
		err := handlers.SpaceHandler.Create(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestSpaceHandler_Put(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newSpcHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	inst, err := samples.CreateInstance(mng)
	if err != nil {
		assert.NoError(t, err)
		return
	}

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating space. Instance not found",
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating space.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating space.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/spaces",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)
		err := handlers.SpaceHandler.Put(ctx)

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

func TestSpaceHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/spaces", tt.Body, tt.Params, nil)
		err := handlers.SpaceHandler.Put(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestSpaceHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newSpcHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	inst, err := samples.CreateInstance(mng)
	if err != nil {
		assert.NoError(t, err)
		return
	}
	space := space.New()
	space.Name = "name"
	space.Desc = "desc"
	space.InstID = inst.ID
	created, _, err := srv.Create(&space)

	if assert.NoError(t, err) {
		assert.True(t, created, "Space created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding space. Ok",
				params: map[string]string{"id": space.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding space. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/spaces/:id", "", tt.params, nil)
			err := handlers.SpaceHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						fmt.Sprintf(
							`"name":"name","description":"desc","instanceId":"%s"`,
							space.InstID),
						strings.Concat(tt.name, "Get space"))
				}
				assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			}
		}
	}
}

func TestSpaceHandler_GetWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/spaces/:id", "", tt.Param, nil)
		err := handlers.SpaceHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestSpaceHandler_ListEntes(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newSpcHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, ente, err := entesmpl.CreateLinkForSpace(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/spaces/:id/entes",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "name"})

		err := handlers.SpaceHandler.ListEntes(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"name","enteId":"%s"}]`, ente.ID),
				"List entes")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestSpaceHandler_GetListEntesError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/spaces/:id/entes", "", tt.Param, tt.QParam)
		err := handlers.SpaceHandler.ListEntes(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestSpaceHandler_ListCategories(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newSpcHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, ente, err := catsmpl.CreateRootLink(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/categories/:id/categories",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "name"})

		err := handlers.SpaceHandler.ListCategories(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"name","categoryId":"%s"}]`, ente.ID),
				"List categories")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestSpaceHandler_GetListCategoriesError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categories/:id/categories", "", tt.Param, tt.QParam)
		err := handlers.SpaceHandler.ListCategories(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newSpcHandlerFaked(t *testing.T) (*runtime.Container, Handlers, space.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	srv := space.NewService(cnt, ext, crud)
	hands := Handlers{SpaceHandler: NewSpaceHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newSpcHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	srv := space.NewService(cnt, ext, crud)
	hands := Handlers{SpaceHandler: NewSpaceHandle(srv, cnt)}
	return cnt, hands, crud
}
