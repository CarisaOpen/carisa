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

	"github.com/carisa/api/internal/entity"

	"github.com/carisa/api/internal/category"

	"github.com/carisa/api/internal/service"

	"github.com/rs/xid"

	csamples "github.com/carisa/api/internal/category/samples"
	tsamples "github.com/carisa/api/internal/samples"
	"github.com/carisa/api/internal/space/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/api/internal/runtime"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestCategoryHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
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
			name:   "Creating category into space.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","parentId":"%s","root":true`, space.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Creating category. Space not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","parentId":"%s","root":true`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/categories", strings.Concat("{", tt.body, "}"), nil, nil)
		err := handlers.CategoryHandler.Create(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var cat category.Category
				errj := json.NewDecoder(rec.Body).Decode(&cat)
				if assert.NoError(t, errj, tt.name) {
					assert.NotEmpty(t, cat.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestCategoryHandler_CreateWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/categories", tt.Body, nil, nil)
		err := handlers.CategoryHandler.Create(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_Put(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
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
			name:   "Creating category. Space not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","parentId":"%s","root":true`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating category.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","parentId":"%s","root":true`, space.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating category.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc","parentId":"%s","root":true`, space.ID.String()),
			status: nethttp.StatusOK,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/categories",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)
		err := handlers.CategoryHandler.Put(ctx)

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

func TestCategoryHandler_PutWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/categories", tt.Body, tt.Params, nil)
		err := handlers.CategoryHandler.Put(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_Get(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	spc, err := samples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}
	cat := category.New()
	cat.Name = "cname"
	cat.Desc = "cdesc"
	cat.ParentID = spc.ID
	cat.Root = true
	created, _, err := srv.Create(&cat)

	if assert.NoError(t, err) {
		assert.True(t, created, "Category created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding category. Ok",
				params: map[string]string{"id": cat.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding category. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categories/:id", "", tt.params, nil)
			err := handlers.CategoryHandler.Get(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						fmt.Sprintf(
							`"name":"cname","description":"cdesc","parentId":"%s","root":true`,
							cat.ParentKey()),
						"Get category")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestCategoryHandler_GetWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categories/:id", "", tt.Param, nil)
		err := handlers.CategoryHandler.Get(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_ListCategories(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, prop, err := csamples.CreateLink(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/categories/:id/child",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "name"})

		err := handlers.CategoryHandler.ListCategories(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"name","linkId":"%s","category":true}]`, prop.Key()),
				"List categories of the space")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestCategoryHandler_GetListCategoriesError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categories/:id/child", "", tt.Param, tt.QParam)
		err := handlers.CategoryHandler.ListCategories(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_ListProps(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	_, prop, err := csamples.CreateLinkProp(mng, xid.NilID())

	if assert.NoError(t, err) {
		rec, ctx := h.NewHTTP(
			nethttp.MethodGet,
			"/api/categories/:id/properties",
			"",
			map[string]string{"id": xid.NilID().String()},
			map[string]string{"sname": "name"})

		err := handlers.CategoryHandler.ListProps(ctx)
		if assert.NoError(t, err) {
			assert.Contains(
				t,
				rec.Body.String(),
				fmt.Sprintf(`[{"name":"namep","categoryPropId":"%s"}]`, prop.Key()),
				"List properties of the category")
			assert.Equal(t, nethttp.StatusOK, rec.Code, "Http status")
		}
	}
}

func TestCategoryHandler_GetListPropsError(t *testing.T) {
	tests := tsamples.TestListError()

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categories/:id/properties", "", tt.Param, tt.QParam)
		err := handlers.CategoryHandler.ListProps(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_CreateProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	cat, err := csamples.CreateCat(mng)
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
			name:   "Creating category property.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","categoryId":"%s"`, cat.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Creating category property. Category not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","categoryId":"%s"`, xid.NilID()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating category property. Category not found.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","type":1,"categoryId":"%s"`, xid.NilID()),
			status: nethttp.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(nethttp.MethodPost, "/api/categoriesprop", strings.Concat("{", tt.body, "}"), nil, nil)
		err := handlers.CategoryHandler.CreateProp(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if assert.NoError(t, err, tt.name) {
			assert.Equal(t, tt.status, rec.Code, strings.Concat(tt.name, "Http status"))
			if rec.Code == nethttp.StatusCreated {
				assert.Contains(t, rec.Body.String(), tt.body, strings.Concat(tt.name, "Created"))
				var prop category.Prop
				errj := json.NewDecoder(rec.Body).Decode(&prop)
				if assert.NoError(t, errj) {
					assert.NotEmpty(t, prop.ID.String(), strings.Concat(tt.name, "ID no empty"))
				}
			}
		}
	}
}

func TestCategoryHandler_CreatePropWithError(t *testing.T) {
	tests := tsamples.TestCreateWithError("CreateWithRel")

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPost, "/api/categoriesprop", tt.Body, nil, nil)
		err := handlers.CategoryHandler.CreateProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_PutProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	cat, err := csamples.CreateCat(mng)
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
			name:   "Creating category property. Category not found",
			body:   fmt.Sprintf(`"name":"name","description":"desc","categoryId":"%s"`, xid.New().String()),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Creating category property.",
			body:   fmt.Sprintf(`"name":"name","description":"desc","categoryId":"%s"`, cat.ID.String()),
			status: nethttp.StatusCreated,
		},
		{
			name:   "Updating category property.",
			body:   fmt.Sprintf(`"name":"name1","description":"desc1","categoryId":"%s"`, cat.ID.String()),
			status: nethttp.StatusOK,
		},
		{
			name:   "Creating category property. Type can not be changed",
			body:   fmt.Sprintf(`"name":"name","description":"desc","categoryId":"%s", "type":1`, xid.New().String()),
			status: nethttp.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		rec, ctx := h.NewHTTP(
			nethttp.MethodPut,
			"/api/categoriesprop",
			strings.Concat("{", tt.body, "}"),
			params,
			nil)
		err := handlers.CategoryHandler.PutProp(ctx)

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

func TestCategoryHandler_PutPropWithError(t *testing.T) {
	params := map[string]string{"id": xid.NilID().String()}
	tests := tsamples.TestPutWithError("PutWithRel", params)

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodPut, "/api/categoriesprop", tt.Body, tt.Params, nil)
		err := handlers.CategoryHandler.PutProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func TestCategoryHandler_GetProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	cat, err := csamples.CreateCat(mng)
	if err != nil {
		assert.Error(t, err)
		return
	}
	prop := category.NewProp()
	prop.Name = "namep"
	prop.Desc = "descp"
	prop.CatID = cat.ID
	prop.Type = entity.Integer
	created, _, err := srv.CreateProp(&prop)

	if assert.NoError(t, err) {
		assert.True(t, created, "Category property created")

		tests := []struct {
			name   string
			params map[string]string
			status int
		}{
			{
				name:   "Finding property of the category. Ok",
				params: map[string]string{"id": prop.ID.String()},
				status: nethttp.StatusOK,
			},
			{
				name:   "Finding property of the category. Not found",
				params: map[string]string{"id": xid.NilID().String()},
				status: nethttp.StatusNotFound,
			},
		}

		for _, tt := range tests {
			rec, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categoriesprop/:id", "", tt.params, nil)
			err := handlers.CategoryHandler.GetProp(ctx)

			if assert.NoError(t, err) {
				if tt.status == nethttp.StatusOK {
					assert.Contains(
						t,
						rec.Body.String(),
						fmt.Sprintf(
							`"name":"namep","description":"descp","categoryId":"%s","type":1`,
							prop.ParentKey()),
						"Get property of the category")
				}
				assert.Equal(t, tt.status, rec.Code, "Http status")
			}
		}
	}
}

func TestCategoryHandler_GetPropWithError(t *testing.T) {
	tests := tsamples.TestGetWithError()

	h := mock.HTTP()
	cnt, handlers, crud := newCategoryHandlerMocked()
	defer h.Close(cnt.Log)

	for _, tt := range tests {
		if tt.MockOper != nil {
			tt.MockOper(crud)
		}
		_, ctx := h.NewHTTP(nethttp.MethodGet, "/api/categoriesprop/:id", "", tt.Param, nil)
		err := handlers.CategoryHandler.GetProp(ctx)

		assert.Equal(t, tt.Status, err.(*echo.HTTPError).Code, tt.Name)
		assert.Error(t, err, tt.Name)
	}
}

func newCategoryHandlerFaked(t *testing.T) (*runtime.Container, Handlers, category.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	entesrv := ente.NewService(cnt, ext, crud)
	srv := category.NewService(cnt, ext, crud, &entesrv)
	hands := Handlers{CategoryHandler: NewCatHandle(srv, cnt)}
	return cnt, hands, srv, mng
}

func newCategoryHandlerMocked() (*runtime.Container, Handlers, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	entesrv := ente.NewService(cnt, ext, crud)
	srv := category.NewService(cnt, ext, crud, &entesrv)
	hands := Handlers{CategoryHandler: NewCatHandle(srv, cnt)}
	return cnt, hands, crud
}
