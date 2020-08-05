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

	"github.com/rs/xid"

	"github.com/carisa/api/internal/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/api/internal/space"

	"github.com/carisa/api/internal/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/mock"

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
		assert.Error(t, err)
		return
	}

	tests := []struct {
		body   string
		status int
	}{
		{
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/spaces", strings.Concat("{", tt.body, "}"), nil)
		err := handlers.SpaceHandler.Create(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err) {
			assert.Equal(t, tt.status, rec.Code, "Http status")
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, "Created")
				var spc space.Space
				errJ := json.NewDecoder(rec.Body).Decode(&spc)
				if assert.NoError(t, errJ) {
					assert.NotEmpty(t, spc.ID.String(), "ID no empty")
				}
			}
		}
	}
}

func TestSpaceHandler_CreateWithError(t *testing.T) {
	tests := samples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/spaces", tt.Body, nil)
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
		assert.Error(t, err)
		return
	}

	params := map[string]string{"id": xid.NilID().String()}

	tests := []struct {
		body   string
		status int
	}{
		{
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			body:   fmt.Sprintf(`"name":"name1","description":"desc","instanceId":"%s"`, inst.ID.String()),
			status: nethttp.StatusOK,
		},
		{
			body:   fmt.Sprintf(`"name":"name","description":"desc","instanceId":"%s"`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/spaces",
			strings.Concat("{", tt.body, "}"),
			params)
		err := handlers.SpaceHandler.Put(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err) {
			assert.Equal(t, tt.status, rec.Code, "Http status")
			if tt.status != nethttp.StatusNotFound {
				assert.Contains(t, rec.Body.String(), tt.body, "Put")
			}
		}
	}
}

func TestSpaceHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := samples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newSpcHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/spaces", tt.Body, tt.Params)
		err := handlers.SpaceHandler.Put(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newSpcHandlerFaked(t *testing.T) (*runtime.Container, Handlers, space.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	srv := space.NewService(cnt, crud)
	hands := Handlers{SpaceHandler: NewSpaceHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newSpcHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	srv := space.NewService(cnt, crud)
	hands := Handlers{SpaceHandler: NewSpaceHandle(srv, cnt)}
	return cnt, hands, crud
}
