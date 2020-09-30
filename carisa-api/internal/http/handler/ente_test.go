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

	csamples "github.com/carisa/api/internal/category/samples"

	"github.com/carisa/api/internal/service"

	"github.com/carisa/api/internal/ente"

	"github.com/rs/xid"

	esamples "github.com/carisa/api/internal/ente/samples"
	tsamples "github.com/carisa/api/internal/samples"
	"github.com/carisa/api/internal/space/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/api/internal/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestEnteHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	space, err := samples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating ente.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, space.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Creating ente. Space not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/entes", strings.Concat("{", tt.body, "}"), nil, nil)
		err := handlers.EnteHandler.Create(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var ente ente.Ente
				errj := json.NewDecoder(rec.Body).Decode(&ente)
				if assert.NoError(t, errj) {
					assert.NotEmpty(t, ente.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestEnteHandler_CreateWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/entes", tt.Body, nil, nil)
		err := handlers.EnteHandler.Create(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_Put(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	space, err := samples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating ente. Space not found",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating ente.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, space.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating ente.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc","spaceId":"%s"`, space.ID.String()),
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/entes",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)
		err := handlers.EnteHandler.Put(ctx)

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

func TestEnteHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/entes", tt.Body, tt.Params, nil)
		err := handlers.EnteHandler.Put(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	spc, err := samples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}
	ente := ente.New()
	ente.Name = "ename"
	ente.Desc = "edesc"
	ente.SpaceID = spc.ID
	created, _, err := srv.Create(&ente)

	if assert.NoError(t, err) {
		assert.True(t, created, "Ente created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding ente. Ok",
				params: map[string]string{"id": ente.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding ente. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entes/:id", "", tt.params, nil)
			err := handlers.EnteHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						fmt.Sprintf(
							`"name":"ename","description":"edesc","spaceId":"%s"`,
							ente.ParentKey()),
						"Get ente")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestEnteHandler_GetWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entes/:id", "", tt.Param, nil)
		err := handlers.EnteHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_ListProps(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, prop, err := esamples.CreateLinkProp(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/entes/:id/properties",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "namep"})

		err := handlers.EnteHandler.ListProps(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"namep","entePropId":"%s"}]`, prop.Key()),
				"List properties of the ente")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestEnteHandler_GetListPropsError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entes/:id/properties", "", tt.Param, tt.QParam)
		err := handlers.EnteHandler.ListProps(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_LinkToCategory(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	e, err := esamples.CreateEnte(mng)
	if err != nil {
		assert.Error(t, err, "Creating ente")
	}
	cat, err := csamples.CreateCat(mng)
	if err != nil {
		assert.Error(t, err, "Creating category")
	}

	test := []struct {
		name   string
		params map[string]string
		status int
	}{
		{
			name: "Ente not found.",
			params: map[string]string{
				"enteId":     xid.New().String(),
				"categoryId": cat.ID.String(),
			},
			status: nethttp.StatusNotFound,
		},
		{
			name: "Category not found.",
			params: map[string]string{
				"enteId":     e.ID.String(),
				"categoryId": xid.New().String(),
			},
			status: nethttp.StatusNotFound,
		},
		{
			name: "Ente connected.",
			params: map[string]string{
				"enteId":     e.ID.String(),
				"categoryId": cat.ID.String(),
			},
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range test {
		rec, ctx := h.NewHTTP(nethttp.MethodPut, "/api/entes/:enteId/linktocategory:categoryId", "", tt.params, nil)
		err := handlers.EnteHandler.LinkToCat(ctx)
		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.status, rec.Code, "Http status")
		}
	}
}

func TestEnteHandler_LinkToCategoryError(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	tests := []struct {
		name     string
		param    map[string]string
		mockOper func(txn *storage.ErrMockCRUDOper)
		status   int
	}{
		{
			name:   "Param ente not found. Bad request",
			param:  map[string]string{"ente": ""},
			status: nethttp.StatusBadRequest,
		},
		{
			name:   "Param category wrong. Bad request",
			param:  map[string]string{"enteId": xid.New().String(), "categoryId": "123"},
			status: nethttp.StatusBadRequest,
		},
		{
			name:     "LinkTo. Internal server error",
			param:    map[string]string{"enteId": xid.New().String(), "categoryId": xid.New().String()},
			mockOper: func(s *storage.ErrMockCRUDOper) { s.Activate("LinkTo") },
			status:   nethttp.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		if tt.mockOper != nil {
			tt.mockOper(crud)
		}

		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entes/:enteId/linktocategory:categoryId", "", tt.param, nil)
		err := handlers.EnteHandler.LinkToCat(ctx)
		assert.Equal(t, tt.status, err.(*echo.HTTPError).Code, tt.name)
		assert.Error(t, err, tt.name)
	}
}

func TestEnteHandler_CreateProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	e, err := esamples.CreateEnte(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating ente property.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","enteId":"%s","type":2`, e.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Creating ente property. Ente not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","enteId":"%s","type":2`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/entesprop", strings.Concat("{", tt.body, "}"), nil, nil)
		err := handlers.EnteHandler.CreateProp(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var prop ente.Prop
				errj := json.NewDecoder(rec.Body).Decode(&prop)
				if assert.NoError(t, errj) {
					assert.NotEmpty(t, prop.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestEnteHandler_CreatePropWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/entesprop", tt.Body, nil, nil)
		err := handlers.EnteHandler.CreateProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_PutProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	ente, err := esamples.CreateEnte(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "Creating ente property. Ente not found",
			body:   fmt.Sprintf(`"name":"name","description":"desc","enteId":"%s","type":2`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating ente property.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","enteId":"%s","type":2`, ente.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating ente property.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc","enteId":"%s","type":3`, ente.ID.String()),
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/entesprop",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)
		err := handlers.EnteHandler.PutProp(ctx)

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

func TestEnteHandler_PutPropWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/entesprop", tt.Body, tt.Params, nil)
		err := handlers.EnteHandler.PutProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestEnteHandler_GetProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newEnteHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	e, err := esamples.CreateEnte(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}
	prop := ente.NewProp()
	prop.Name = "namep"
	prop.Desc = "descp"
	prop.EnteID = e.ID
	created, _, err := srv.CreateProp(&prop)

	if assert.NoError(t, err) {
		assert.True(t, created, "Ente property created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding property of the ente. Ok",
				params: map[string]string{"id": prop.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding property of the ente. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entesprop/:id", "", tt.params, nil)
			err := handlers.EnteHandler.GetProp(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						fmt.Sprintf(
							`"name":"namep","description":"descp","enteId":"%s","type":1`,
							prop.ParentKey()),
						"Get property of the ente")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestEnteHandler_GetPropWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newEnteHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/entesprop/:id", "", tt.Param, nil)
		err := handlers.EnteHandler.GetProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newEnteHandlerFaked(t *testing.T) (*runtime.Container, Handlers, ente.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	srv := ente.NewService(cnt, ext, crud)
	hands := Handlers{EnteHandler: NewEnteHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newEnteHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	srv := ente.NewService(cnt, ext, crud)
	hands := Handlers{EnteHandler: NewEnteHandle(srv, cnt)}
	return cnt, hands, crud
}
