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
	"context"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"testing"

	"github.com/carisa/internal/api/category"
	"github.com/carisa/internal/api/ente"
	esamples "github.com/carisa/internal/api/ente/samples"
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"

	"github.com/rs/xid"

	csamples "github.com/carisa/internal/api/category/samples"
	tsamples "github.com/carisa/internal/api/samples"
	"github.com/carisa/internal/api/space/samples"

	"github.com/labstack/echo/v4"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

func TestCategoryHandler_Create(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, srv, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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

func TestCategoryHandler_LinkToProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, srv, srve, mng := newCategoryHandlerFaked(t)
	defer mng.Close()
	defer h.Close(cnt.Log)

	catRoot, err := csamples.CreateCat(mng)
	if err != nil {
		assert.Error(t, err, "Creating root category")
		return
	}
	catPropRoot := category.NewProp()
	catPropRoot.CatID = catRoot.ID
	_, _, err = srv.CreateProp(&catPropRoot)
	if err != nil {
		assert.Error(t, err, "Creating root category property")
		return
	}

	ok, catChild, catChildProp1 := createCat(t, srv, catRoot, entity.Integer)
	if !ok {
		return
	}
	catChildProp2 := category.NewProp()
	catChildProp2.CatID = catChild.ID
	catChildProp2.Type = entity.Boolean
	_, _, err = srv.CreateProp(&catChildProp2)
	if err != nil {
		assert.Error(t, err, "Creating a second property in the child category")
		return
	}

	ok, _, catccProp := createCat(t, srv, catChild, entity.Integer)
	if !ok {
		return
	}

	enteChild, err := esamples.CreateEnte(mng)
	if err != nil {
		assert.Error(t, err, "Creating child ente")
		return
	}
	enteChildProp := ente.NewProp()
	enteChildProp.Type = entity.Integer
	enteChildProp.EnteID = enteChild.ID
	enteChildProp.Name = "nameep"
	_, _, err = srve.CreateProp(&enteChildProp)
	if err != nil {
		assert.Error(t, err, "Creating child ente property")
		return
	}
	_, _, _, err = srve.LinkToCat(enteChild.ID, catRoot.ID)
	if err != nil {
		assert.Error(t, err, "Creating linking between category root and ente")
		return
	}

	tests := []struct {
		name    string
		source  xid.ID
		target  xid.ID
		status  int
		resBody string
		typep   entity.TypeProp
	}{
		{
			name:   "Category property not found",
			source: xid.NilID(),
			target: xid.NilID(),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "Target property not found",
			source: catPropRoot.ID,
			target: xid.NilID(),
			status: nethttp.StatusNotFound,
		},
		{
			name:   "The category or ente of the property is not child of the category of the property",
			source: catPropRoot.ID,
			target: catccProp.ID,
			status: nethttp.StatusBadRequest,
		},
		{
			name:    "The category property is linked successfully with other category property",
			source:  catPropRoot.ID,
			target:  catChildProp1.ID,
			status:  nethttp.StatusOK,
			resBody: fmt.Sprintf(`{"name":"namecp","propertyId":"%s","category":true}`, catChildProp1.ID.String()),
			typep:   entity.Integer,
		},
		{
			name:   "The category property is not the same type than the target property",
			source: catPropRoot.ID,
			target: catChildProp2.ID,
			status: nethttp.StatusConflict,
		},
		{
			name:    "The category property is linked successfully with a ente property",
			source:  catPropRoot.ID,
			target:  enteChildProp.ID,
			status:  nethttp.StatusOK,
			resBody: fmt.Sprintf(`{"name":"nameep","propertyId":"%s","category":false}`, enteChildProp.ID.String()),
			typep:   entity.Integer,
		},
	}

	for _, tt := range tests {
		params := map[string]string{"catPropId": tt.source.String(), "propId": tt.target.String()}
		rec, ctx := h.NewHTTP(nethttp.MethodPut, "/api/categoriesProp/:catPropId/linkto/:propId", "", params, nil)
		err := handlers.CategoryHandler.LinkToProp(ctx)

		if err != nil && tt.status == err.(*echo.HTTPError).Code {
			continue
		}
		if err != nil {
			assert.Error(t, err)
			return
		}

		assert.Contains(t, rec.Body.String(), tt.resBody, tt.name)
		var catp category.Prop
		_, err = srv.GetProp(catPropRoot.ID, &catp)
		if assert.NoError(t, err) {
			assert.Equal(t, tt.typep, catp.Type, tt.name)
		}
		found, err := mng.Store().Exists(context.TODO(), storage.DLRKey(tt.target.String(), tt.source.String()))
		if assert.NoError(t, err) {
			assert.True(t, found, "Getting DLR")
		}
	}
}

func createCat(t *testing.T,
	service category.Service,
	catParent category.Category,
	typep entity.TypeProp) (bool, category.Category, category.Prop) {
	//
	catChild := category.New()
	catChild.ParentID = catParent.ID
	_, _, err := service.Create(&catChild)
	if err != nil {
		assert.Error(t, err, "Creating child category")
		return false, category.Category{}, category.Prop{}
	}
	catChildProp := category.NewProp()
	catChildProp.CatID = catChild.ID
	catChildProp.Name = "namecp"
	catChildProp.Type = typep
	_, _, err = service.CreateProp(&catChildProp)
	if err != nil {
		assert.Error(t, err, "Creating child category property")
		return false, category.Category{}, category.Prop{}
	}
	return true, catChild, catChildProp
}

func TestCategoryHandler_CreateProp(t *testing.T) {
	h := mock.HTTP()
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, _, _, mng := newCategoryHandlerFaked(t)
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
	cnt, handlers, srv, _, mng := newCategoryHandlerFaked(t)
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

func newCategoryHandlerFaked(t *testing.T) (*runtime.Container, Handlers, category.Service, ente.Service, storage.Integration) {
	mng, cnt, crud := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crud.Store())
	entesrv := ente.NewService(cnt, ext, crud)
	srv := category.NewService(cnt, ext, crud, &entesrv)
	hands := Handlers{CategoryHandler: NewCatHandle(srv, cnt)}
	return cnt, hands, srv, entesrv, mng
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
