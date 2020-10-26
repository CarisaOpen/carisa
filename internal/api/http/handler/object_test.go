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
	"fmt"
	nethttp "net/http"
	"testing"

	osamples "github.com/carisa/internal/api/object/samples"

	"github.com/carisa/internal/api/samples"

	psamples "github.com/carisa/internal/api/plugin/samples"
	tsamples "github.com/carisa/internal/api/samples"

	"github.com/carisa/internal/api/object"
	"github.com/carisa/internal/api/plugin"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/mock"
	"github.com/carisa/pkg/strings"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestObjectHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newObjectHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	container, err := samples.CreateEntityMock(mng)
	if err != nil {
		assert.Error(t, err, "Creating container")
		return
	}

	plugins := plugin.Plugins()
	protoID := xid.New()

	tests := []struct {
		name   string
		body   string
		status int
		errmsg string
	}{
		{
			name: "Creating object instance.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				container.ID.String(),
				protoID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name: "Creating object instance. Prototype not found.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				container.ID.String(),
				xid.NilID()),
			status: nethttp.StatusNotFound,
			errmsg: "code=404, message=[the plugin prototype not found]",
		},
		{
			name: "Creating object instance. Container not found.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				xid.NilID(),
				protoID.String()),
			status: nethttp.StatusNotFound,
			errmsg: "code=404, message=[container not found]",
		},
	}

	for _, pc := range plugins {
		_, err := psamples.CreatePlugin(mng, pc, protoID)
		if err != nil {
			assert.Error(t, err, "Creating plugin")
			continue
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/objects", strings.Concat("{", tt.body, "}"), nil, nil)
			err := handlers.ObjectHandler.Create(ctx, pc)

			if err != nil && tt.status == err.(*echo.HTTPError).Code {
				assert.Equal(t, tt.errmsg, err.Error(), strings.Concat(tt.name, "Error message"))
				continue
			}
			if err != nil {
				assert.Error(t, err, "Error creating")
			}
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, strings.Concat(tt.name, "Http status")))
			if rec.Code != nethttp.StatusCreated {
				continue
			}
			assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
		}
	}
}

func TestObjectHandler_CreateWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	plugins := plugin.Plugins()

	h := mock.HTTP()
	cnt, handlers, crud := newObjectHandlerMocked()
	defer h.Close(cnt.Log)

	for _, pc := range plugins {
		for _, tt := range tests {
			if tt.MockOper != nil {
				tt.MockOper(crud)
			}
			_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/objects", tt.Body, nil, nil)
			err := handlers.ObjectHandler.Create(ctx, pc)

			assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
			assert.Error(t, err, tt.Name)
		}
	}
}

func TestObjectHandler_Put(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newObjectHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	container, err := samples.CreateEntityMock(mng)
	if err != nil {
		assert.Error(t, err, "Creating container")
		return
	}

	plugins := plugin.Plugins()
	protoID := xid.New()

	tests := []struct {
		name   string
		body   string
		status int
		errmsg string
	}{
		{
			name: "Creating object instance. Prototype not found.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				container.ID.String(),
				xid.NilID()),
			status: nethttp.StatusNotFound,
			errmsg: "code=404, message=[the plugin prototype not found]",
		},
		{
			name: "Creating object instance. Container not found.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				xid.NilID(),
				protoID.String()),
			status: nethttp.StatusNotFound,
			errmsg: "code=404, message=[container not found]",
		},
		{
			name: "Creating object instance.",
			body: fmt.Sprintf(
				`"name":"name","description":"desc","containerId":"%s","prototypeId":"%s"`,
				container.ID.String(),
				protoID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name: "Updating object instance.",
			body: fmt.Sprintf(
				`"name":"name1","description":"desc1","containerId":"%s","prototypeId":"%s"`,
				container.ID.String(),
				protoID.String()),
			status: nethttp.StatusOK,
		},
	}

	for _, pc := range plugins {
		_, err := psamples.CreatePlugin(mng, pc, protoID)
		if err != nil {
			assert.Error(t, err, "Creating plugin")
			return
		}
		params := map[string]string{"id": xid.NilID().String()}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(
				nethttp.MethodPut,
				"/api/objects",
				strings.Concat("{", tt.body, "}"),
				params,
				nil)
			err := handlers.ObjectHandler.Put(ctx, pc)

			if err != nil && tt.status == err.(*echo.HTTPError).Code {
				assert.Equal(t, tt.errmsg, err.Error(), strings.Concat(tt.name, "Error message"))
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

func TestObjectHandler_PutWithError(t *testing.T) {
	cnt, handlers, crud := newObjectHandlerMocked()
	h := mock.HTTP()
	defer h.Close(cnt.Log)

	plugins := plugin.Plugins()

	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	for _, pc := range plugins {
		for _, tt := range tests {
			if tt.MockOper != nil {
				tt.MockOper(crud)
			}
			_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/objects", tt.Body, tt.Params, nil)
			err := handlers.ObjectHandler.Put(ctx, pc)

			assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
			assert.Error(t, err, tt.Name)
		}
	}
}
func TestObjectHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newObjectHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	container, err := samples.CreateEntityMock(mng)
	if err != nil {
		assert.Error(t, err, "Creating container")
		return
	}

	instID := xid.New()

	tests := []struct {
		name   string
		params map[string]string
		status int
	}{
		{
			name:   "Finding inst. Ok",
			params: map[string]string{"id": instID.String()},
			status: nethttp.StatusOK,
		},
		{
			name:   "Finding inst. Not found",
			params: map[string]string{"id": xid.NilID().String()},
			status: nethttp.StatusNotFound,
		},
	}

	plugins := plugin.Plugins()

	for _, pc := range plugins {
		protoID := xid.New()
		_, err := psamples.CreatePlugin(mng, pc, protoID)
		if err != nil {
			assert.Error(t, err, "Creating plugin")
			continue
		}

		inst := object.New()
		inst.ID = instID
		inst.Name = "iname"
		inst.Desc = "idesc"
		inst.ContainerID = container.ID
		inst.ProtoID = protoID
		_, _, _, err = srv.Put(&inst)
		if err != nil {
			assert.Error(t, err, "Creating instance")
			return
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/objects", "", tt.params, nil)
			err := handlers.ObjectHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						`"name":"iname","description":"idesc"`,
						"Get instance")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestObjectHandler_GetWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newObjectHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/objects", "", tt.Param, nil)
		err := handlers.ObjectHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestObjectHandler_ListQueries(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newObjectHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	plugins := plugin.Plugins()

	for _, pc := range plugins {
		_, prop, err := osamples.CreateLink(mng, xid.NilID(), pc)

		if assert.NoError(t, err) {
			rec, ctx := h.NewHTTP(
				nethttp.MethodGet,
				"/api/queries",
				"",
				map[string]string{"id": xid.NilID().String()},
				map[string]string{"sname": "name"})

			err := handlers.ObjectHandler.ListInstances(ctx, pc)
			if assert.NoError(t, err) {
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`[{"name":"name","instanceId":"%s","category":"%s"}]`, prop.Key(), string(pc)),
					"List categories of the space")
				assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
			}
		}
	}
}

func newObjectHandlerFaked(t *testing.T) (*runtime.Container, Handlers, object.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	plugin := plugin.NewService(cnt, ext, crud)
	srv := object.NewService(cnt, ext, crud, &plugin)
	hands := Handlers{ObjectHandler: NewObjectHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newObjectHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	plugin := plugin.NewService(cnt, ext, crud)
	srv := object.NewService(cnt, ext, crud, &plugin)
	hands := Handlers{ObjectHandler: NewObjectHandle(srv, cnt)}
	return cnt, hands, crud
}
