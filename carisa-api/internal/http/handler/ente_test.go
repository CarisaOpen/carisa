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

	"github.com/carisa/api/internal/ente"

	"github.com/rs/xid"

	tsamples "github.com/carisa/api/internal/samples"
	"github.com/carisa/api/internal/space/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/api/internal/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestEnteService_Create(t *testing.T) {
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
			name:   "Creating ente. Instance not found.",
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
		if assert.NoError(t, err) {
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

func TestEnteService_CreateWithError(t *testing.T) {
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

func TestEnteService_Put(t *testing.T) {
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
			name:   "Creating ente.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, space.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating ente.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc","spaceId":"%s"`, space.ID.String()),
			status: nethttp.StatusOK,
		},
		{
			name:   "Creating ente. Space not found",
			body:   fmt.Sprintf(`"name":"name","description":"desc","spaceId":"%s"`, xid.New().String()),
			status: nethttp.StatusNotFound,
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

func TestEnteService_PutWithError(t *testing.T) {
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

func TestEnteService_Get(t *testing.T) {
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

func TestEnteService_GetWithError(t *testing.T) {
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

func newEnteHandlerFaked(t *testing.T) (*runtime.Container, Handlers, ente.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	srv := ente.NewService(cnt, crud)
	hands := Handlers{EnteHandler: NewEnteHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newEnteHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	srv := ente.NewService(cnt, crud)
	hands := Handlers{EnteHandler: NewEnteHandle(srv, cnt)}
	return cnt, hands, crud
}
