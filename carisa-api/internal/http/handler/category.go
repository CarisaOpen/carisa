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
	if err = errService(
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
	if err = errService(
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
