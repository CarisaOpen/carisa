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

	"github.com/carisa/api/internal/space"

	"github.com/carisa/pkg/http"

	"github.com/carisa/api/internal/http/validator"

	"github.com/carisa/api/internal/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locSpace = "http.space"

// Space hands the http request of the instance
type Space struct {
	srv space.Service
	cnt *runtime.Container
}

// NewSpaceHandle creates handler
func NewSpaceHandle(srv space.Service, cnt *runtime.Container) Space {
	return Space{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the space domain
func (s *Space) Create(c httpc.Context) error {
	spc := space.Space{}
	if err := c.Bind(&spc); err != nil {
		return s.ErrorRecover(c, err)
	}

	if httpErr := validator.Descriptor(c, spc.Descriptor); httpErr != nil {
		return httpErr
	}

	created, found, err := s.srv.Create(&spc)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create the space")
	}
	if !found {
		return c.HTTPError(nethttp.StatusNotFound, "instance not found")
	}

	return c.JSON(http.CreateStatus(created), spc)
}

// Put creates or update the space domain
func (s *Space) Put(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	space := space.Space{}
	if err := c.Bind(&space); err != nil {
		return s.ErrorRecover(c, err)
	}

	if httpErr := validator.Descriptor(c, space.Descriptor); httpErr != nil {
		return httpErr
	}

	space.ID = id
	updated, found, err := s.srv.Put(&space)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create or update the space")
	}
	if !found {
		return c.HTTPError(nethttp.StatusNotFound, "instance not found")
	}

	return c.JSON(http.PutStatus(updated), space)
}

func (s *Space) ErrorRecover(c httpc.Context, err error) error {
	return c.HTTPErrorLog(
		nethttp.StatusBadRequest,
		"cannot recover the space",
		err,
		s.cnt.Log,
		locSpace)
}
