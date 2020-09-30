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
	nethttp "net/http"

	"github.com/carisa/api/internal/entity"

	"github.com/carisa/api/internal/category"

	"github.com/carisa/api/internal/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/api/internal/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locCat = "http.category"

// Category hands the http request of the category
type Category struct {
	srv category.Service
	cnt *runtime.Container
}

// NewCatHandle creates handler
func NewCatHandle(srv category.Service, cnt *runtime.Container) Category {
	return Category{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the category
func (c *Category) Create(ctx httpc.Context) error {
	cat := category.Category{}
	if err := bind(ctx, locCat, c.cnt.Log, &cat); err != nil {
		return err
	}

	created, found, err := c.srv.Create(&cat)
	if err = errCRUDSrv(
		ctx,
		err,
		"it was impossible to create the category",
		"category or space not found", found); err != nil {
		return err
	}

	return ctx.JSON(http.CreateStatus(created), cat)
}

// Put creates or update the category
func (c *Category) Put(ctx httpc.Context) error {
	id, err := convert.ParamID(ctx)
	if err != nil {
		return err
	}

	cat := category.Category{}
	if err := bind(ctx, locCat, c.cnt.Log, &cat); err != nil {
		return err
	}

	cat.ID = id
	updated, found, err := c.srv.Put(&cat)
	if err = errCRUDSrv(
		ctx,
		err,
		"it was impossible to create or update the category",
		"space or category not found", found); err != nil {
		return err
	}

	return ctx.JSON(http.PutStatus(updated), cat)
}

// Get gets the category by ID
func (c *Category) Get(ctx httpc.Context) error {
	var cat category.Category

	id, err := convert.ParamID(ctx)
	if err != nil {
		return err
	}

	found, err := c.srv.Get(id, &cat)
	if err != nil {
		return ctx.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the category")
	}

	return ctx.JSON(http.GetStatus(found), cat)
}

// ListCategories list child categories by category ID and return top categories.
// If sname query param is not empty, is filtered by categories which name starts by name parameter
// If gtname query param is not empty, is filtered by categories which name is greater than name parameter
func (c *Category) ListCategories(ctx httpc.Context) error {
	id, name, top, ranges, err := convert.FilterLink(ctx)
	if err != nil {
		return err
	}

	props, err := c.srv.ListCategories(id, name, ranges, top)
	if err != nil {
		return ctx.HTTPError(nethttp.StatusInternalServerError, "it was impossible to list the child categories of the category")
	}

	return ctx.JSON(nethttp.StatusOK, props)
}

// ListProps list properties by category ID and return top properties.
// If sname query param is not empty, is filtered by properties which name starts by name parameter
// If gtname query param is not empty, is filtered by properties which name is greater than name parameter
func (c *Category) ListProps(ctx httpc.Context) error {
	id, name, top, ranges, err := convert.FilterLink(ctx)
	if err != nil {
		return err
	}

	props, err := c.srv.ListProps(id, name, ranges, top)
	if err != nil {
		return ctx.HTTPError(nethttp.StatusInternalServerError, "it was impossible to list the properties of the category")
	}

	return ctx.JSON(nethttp.StatusOK, props)
}

// CreateProp creates the property of the category
func (c *Category) CreateProp(ctx httpc.Context) error {
	prop := category.Prop{}
	if err := bind(ctx, locCat, c.cnt.Log, &prop); err != nil {
		return err
	}
	if err := validType(ctx, prop); err != nil {
		return err
	}

	created, found, err := c.srv.CreateProp(&prop)
	if err = errCRUDSrv(
		ctx,
		err,
		"it was impossible to create the property of the category",
		"category not found",
		found); err != nil {
		return err
	}

	return ctx.JSON(http.CreateStatus(created), prop)
}

// PutProp creates or update the property of the category
func (c *Category) PutProp(ctx httpc.Context) error {
	id, err := convert.ParamID(ctx)
	if err != nil {
		return err
	}

	prop := category.Prop{}
	if err := bind(ctx, locCat, c.cnt.Log, &prop); err != nil {
		return err
	}
	if err := validType(ctx, prop); err != nil {
		return err
	}

	prop.ID = id
	updated, found, err := c.srv.PutProp(&prop)
	if err = errCRUDSrv(
		ctx, err, "it was impossible to create or update the property of the category", "category not found", found); err != nil {
		return err
	}

	return ctx.JSON(http.PutStatus(updated), prop)
}

// GetProp gets the property of the category by ID
func (c *Category) GetProp(ctx httpc.Context) error {
	var prop category.Prop

	id, err := convert.ParamID(ctx)
	if err != nil {
		return err
	}

	found, err := c.srv.GetProp(id, &prop)
	if err != nil {
		return ctx.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the property of the category")
	}

	return ctx.JSON(http.GetStatus(found), prop)
}

// LinkToProp connects a category property or with other or with a ente property
func (c *Category) LinkToProp(ctx httpc.Context) error {
	catPropId, err := convert.ParamXID(ctx, "catPropId")
	if err != nil {
		return err
	}
	propId, err := convert.ParamXID(ctx, "propId")
	if err != nil {
		return err
	}

	pfound, cfound, isChild, equalType, rel, err := c.srv.LinkToProp(catPropId, propId)
	if err != nil {
		return ctx.HTTPError(nethttp.StatusInternalServerError, err)
	}

	if !pfound {
		return ctx.HTTPError(nethttp.StatusNotFound, "Category property not found")
	}
	if !cfound {
		return ctx.HTTPError(nethttp.StatusNotFound, "Target property (category or ente property) not found")
	}
	if !isChild {
		return ctx.HTTPError(
			nethttp.StatusBadRequest,
			"The category of the target property (category or ente property) must be child of the category of the property for linking")
	}
	if !equalType {
		return ctx.HTTPError(
			nethttp.StatusConflict,
			"The target property (category or ente property) must be of the same type than the category property for linking")
	}

	return ctx.JSON(nethttp.StatusOK, rel)
}

func validType(ctx httpc.Context, prop category.Prop) error {
	if prop.Type != entity.None {
		return ctx.HTTPError(nethttp.StatusBadRequest, "The 'type' property can be changed")
	}
	return nil
}
