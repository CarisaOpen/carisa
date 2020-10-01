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

	"github.com/carisa/internal/api/ente"

	"github.com/carisa/internal/api/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/internal/api/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locEnte = "http.ente"

// Ente hands the http request of the ente
type Ente struct {
	srv ente.Service
	cnt *runtime.Container
}

// NewEnteHandle creates handler
func NewEnteHandle(srv ente.Service, cnt *runtime.Container) Ente {
	return Ente{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the ente
func (e *Ente) Create(c httpc.Context) error {
	ente := ente.Ente{}
	if err := bind(c, locEnte, e.cnt.Log, &ente); err != nil {
		return err
	}

	created, found, err := e.srv.Create(&ente)
	if err = errCRUDSrv(c, err, "it was impossible to create the ente", "space not found", found); err != nil {
		return err
	}

	return c.JSON(http.CreateStatus(created), ente)
}

// Put creates or update the ente
func (e *Ente) Put(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	ente := ente.Ente{}
	if err := bind(c, locEnte, e.cnt.Log, &ente); err != nil {
		return err
	}

	ente.ID = id
	updated, found, err := e.srv.Put(&ente)
	if err = errCRUDSrv(
		c, err, "it was impossible to create or update the ente", "space not found", found); err != nil {
		return err
	}

	return c.JSON(http.PutStatus(updated), ente)
}

// Get gets the ente by ID
func (e *Ente) Get(c httpc.Context) error {
	var ente ente.Ente

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := e.srv.Get(id, &ente)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the ente")
	}

	return c.JSON(http.GetStatus(found), ente)
}

// LinkToCat connects ente to category in the tree
func (e *Ente) LinkToCat(c httpc.Context) error {
	enteID, err := convert.ParamXID(c, "enteId")
	if err != nil {
		return err
	}
	catID, err := convert.ParamXID(c, "categoryId")
	if err != nil {
		return err
	}

	efound, cfound, rel, err := e.srv.LinkToCat(enteID, catID)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, err)
	}

	if !efound {
		return c.HTTPError(nethttp.StatusNotFound, "Ente not found")
	}
	if !cfound {
		return c.HTTPError(nethttp.StatusNotFound, "Category not found")
	}

	return c.JSON(nethttp.StatusOK, rel)
}

// ListProps list properties by ente ID and return top properties.
// If sname query param is not empty, is filtered by properties which name starts by name parameter
// If gtname query param is not empty, is filtered by properties which name is greater than name parameter
func (e *Ente) ListProps(c httpc.Context) error {
	id, name, top, ranges, err := convert.FilterLink(c)
	if err != nil {
		return err
	}

	props, err := e.srv.ListProps(id, name, ranges, top)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to list the properties of the ente")
	}

	return c.JSON(nethttp.StatusOK, props)
}

// CreateProp creates the property of ente
func (e *Ente) CreateProp(c httpc.Context) error {
	prop := ente.Prop{}
	if err := bind(c, locEnte, e.cnt.Log, &prop); err != nil {
		return err
	}

	created, found, err := e.srv.CreateProp(&prop)
	if err = errCRUDSrv(c, err, "it was impossible to create the property of the ente", "ente not found", found); err != nil {
		return err
	}

	return c.JSON(http.CreateStatus(created), prop)
}

// PutProp creates or update the property of ente category
func (e *Ente) PutProp(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	prop := ente.Prop{}
	if err := bind(c, locEnte, e.cnt.Log, &prop); err != nil {
		return err
	}

	prop.ID = id
	updated, found, err := e.srv.PutProp(&prop)
	if err = errCRUDSrv(
		c, err, "it was impossible to create or update the property of the ente", "ente not found", found); err != nil {
		return err
	}

	return c.JSON(http.PutStatus(updated), prop)
}

// GetProp gets the property of ente by ID
func (e *Ente) GetProp(c httpc.Context) error {
	var prop ente.Prop

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := e.srv.GetProp(id, &prop)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the property of the ente")
	}

	return c.JSON(http.GetStatus(found), prop)
}
