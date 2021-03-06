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

	"github.com/carisa/internal/api/http/convert"

	"github.com/carisa/internal/api/space"

	"github.com/carisa/pkg/http"

	"github.com/carisa/internal/api/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locSpace = "http.space"

// Space hands the http request of the space.Space
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

// Create creates the space.Space
func (s *Space) Create(c httpc.Context) error {
	spc := space.Space{}
	if err := bind(c, locSpace, s.cnt.Log, &spc); err != nil {
		return err
	}

	created, found, err := s.srv.Create(&spc)
	if err := errCRUDSrv(
		c, err, "it was impossible to create the space", "instance not found", found); err != nil {
		return err
	}

	return c.JSON(http.CreateStatus(created), spc)
}

// Put creates or update the space.Space
func (s *Space) Put(c httpc.Context) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	spc := space.Space{}
	if err := bind(c, locSpace, s.cnt.Log, &spc); err != nil {
		return err
	}

	spc.ID = id
	updated, found, err := s.srv.Put(&spc)
	if err := errCRUDSrv(
		c, err, "it was impossible to create or update the space", "instance not found", found); err != nil {
		return err
	}

	return c.JSON(http.PutStatus(updated), spc)
}

// Get gets the space.Space by ID
func (s *Space) Get(c httpc.Context) error {
	var space space.Space

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := s.srv.Get(id, &space)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the space")
	}

	return c.JSON(http.GetStatus(found), space)
}

// ListEntes list entes by space.Space ID and return top entes.
// If sname query param is not empty, is filtered by entes which name starts by name parameter
// If gtname query param is not empty, is filtered by entes which name is greater than name parameter
func (s *Space) ListEntes(c httpc.Context) error {
	id, name, top, ranges, err := convert.FilterLink(c, false)
	if err != nil {
		return err
	}

	entes, err := s.srv.ListEntes(id, name, ranges, top)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to list the entes")
	}

	return c.JSON(nethttp.StatusOK, entes)
}

// ListCategories list categories by space.Space ID and return top categories.
// If sname query param is not empty, is filtered by categories which name starts by name parameter
// If gtname query param is not empty, is filtered by categories which name is greater than name parameter
func (s *Space) ListCategories(c httpc.Context) error {
	id, name, top, ranges, err := convert.FilterLink(c, false)
	if err != nil {
		return err
	}

	categories, err := s.srv.ListCategories(id, name, ranges, top)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to list the categories")
	}

	return c.JSON(nethttp.StatusOK, categories)
}
