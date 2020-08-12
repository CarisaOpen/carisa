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

	"github.com/carisa/api/internal/ente"

	"github.com/carisa/api/internal/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/api/internal/runtime"
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

// Create creates the ente domain
func (e *Ente) Create(c httpc.Context) error {
	ente := ente.Ente{}
	if err := bind(c, locEnte, e.cnt.Log, &ente, ente.Descriptor); err != nil {
		return err
	}

	created, found, err := e.srv.Create(&ente)
	if err = errService(c, err, "it was impossible to create the ente", "space not found", found); err != nil {
		return err
	}

	return c.JSON(http.CreateStatus(created), ente)
}

// Put creates or update the ente domain
func (e *Ente) Put(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	ente := ente.Ente{}
	if err := bind(c, locEnte, e.cnt.Log, &ente, ente.Descriptor); err != nil {
		return err
	}

	ente.ID = id
	updated, found, err := e.srv.Put(&ente)
	if err = errService(
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
