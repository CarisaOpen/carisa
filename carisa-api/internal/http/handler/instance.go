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

	"github.com/carisa/api/internal/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/api/internal/http/validator"

	"github.com/carisa/api/internal/instance"
	"github.com/carisa/api/internal/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locInstance = "http.instance"

// Instance hands the http request of the instance
type Instance struct {
	srv instance.Service
	cnt *runtime.Container
}

// NewInstanceHandle creates handler
func NewInstanceHandle(srv instance.Service, cnt *runtime.Container) Instance {
	return Instance{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the instance domain
func (i *Instance) Create(c httpc.Context) error {
	inst := instance.Instance{}
	if err := c.Bind(&inst); err != nil {
		return i.ErrorRecover(c, err)
	}

	if httpErr := validator.Descriptor(c, inst.Descriptor); httpErr != nil {
		return httpErr
	}

	created, err := i.srv.Create(&inst)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create the instance")
	}

	return c.JSON(http.CreateStatus(created), inst)
}

// Put creates or update the instance domain
func (i *Instance) Put(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	inst := instance.Instance{}
	if err := c.Bind(&inst); err != nil {
		return i.ErrorRecover(c, err)
	}

	if httpErr := validator.Descriptor(c, inst.Descriptor); httpErr != nil {
		return httpErr
	}

	inst.ID = id
	updated, err := i.srv.Put(&inst)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create or update the instance")
	}

	return c.JSON(http.PutStatus(updated), inst)
}

func (i *Instance) Get(c httpc.Context) error {
	var inst instance.Instance

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := i.srv.Get(id, &inst)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the instance")
	}

	return c.JSON(http.GetStatus(found), inst)
}

func (i *Instance) ErrorRecover(c httpc.Context, err error) error {
	return c.HTTPErrorLog(
		nethttp.StatusBadRequest,
		"cannot recover the instance",
		err,
		i.cnt.Log,
		locInstance)
}
