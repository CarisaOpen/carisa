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

	"github.com/rs/xid"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/api/object"

	"github.com/carisa/internal/api/plugin"

	"github.com/carisa/internal/api/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/internal/api/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locObject = "http.object"

// Object hands the http request of the object instance
type Object struct {
	srv object.Service
	cnt *runtime.Container
}

// NewObjectHandle creates handler
func NewObjectHandle(srv object.Service, cnt *runtime.Container) Object {
	return Object{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the object instance
func (o *Object) Create(c httpc.Context, cntParam string, schContainer string, category plugin.Category) error {
	cntID, err := convert.ParamXID(c, cntParam)
	if err != nil {
		return err
	}

	inst, err := o.bindInstance(c, category, cntID, schContainer)
	if err != nil {
		return err
	}

	created, foundp, foundc, err := o.srv.Create(&inst)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create the instance")
	}
	if !foundp {
		return c.HTTPError(nethttp.StatusNotFound, "the plugin prototype not found")
	}
	if !foundc {
		return c.HTTPError(nethttp.StatusNotFound, "container not found")
	}
	return c.JSON(http.CreateStatus(created), inst)
}

// Put creates or update the object instance
func (o *Object) Put(c httpc.Context, cntParam string, schContainer string, category plugin.Category) error {
	cntID, err := convert.ParamXID(c, cntParam)
	if err != nil {
		return err
	}
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	inst, err := o.bindInstance(c, category, cntID, schContainer)
	if err != nil {
		return err
	}
	inst.ID = id

	updated, foundp, foundc, err := o.srv.Put(&inst)
	if !foundp {
		return c.HTTPError(nethttp.StatusNotFound, "the plugin prototype not found")
	}
	if err = errCRUDSrv(
		c, err, "it was impossible to create or update the object instance", "container not found", foundc); err != nil {
		return err
	}

	return c.JSON(http.PutStatus(updated), inst)
}

func (o *Object) bindInstance(
	c httpc.Context,
	category plugin.Category,
	cntID xid.ID,
	schContainer string) (object.Instance, error) {
	//
	inst := object.Instance{}
	if err := bind(c, locObject, o.cnt.Log, &inst); err != nil {
		return object.Instance{}, err
	}
	inst.Category = category
	inst.ContainerID = cntID
	inst.SchContainer = schContainer
	return inst, nil
}

// Get gets the object instance by ID
func (o *Object) Get(c httpc.Context) error {
	var inst object.Instance

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := o.srv.Get(id, &inst)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the object instance")
	}

	return c.JSON(http.GetStatus(found), inst)
}

// ListInstances list child queries by ID and return top queries.
// If sname query param is not empty, is filtered by categories which name starts by name parameter
// If gtname query param is not empty, is filtered by categories which name is greater than name parameter
func (o *Object) ListInstances(ctx httpc.Context, schContainer string, category plugin.Category) error {
	id, name, top, ranges, err := convert.FilterLink(ctx, false)
	if err != nil {
		return err
	}

	props, err := o.srv.ListInstances(schContainer, id, category, name, ranges, top)
	if err != nil {
		return ctx.HTTPError(
			nethttp.StatusInternalServerError,
			strings.Concat("it was impossible to list the child ", string(category)))
	}

	return ctx.JSON(nethttp.StatusOK, props)
}
