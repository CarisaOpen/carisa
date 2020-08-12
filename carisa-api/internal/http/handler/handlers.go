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
	echoc "github.com/carisa/pkg/http/echo"
	"github.com/labstack/echo/v4"
)

// Handlers is a handlers container
type Handlers struct {
	InstHandler  Instance
	SpaceHandler Space
	EnteHandler  Ente
}

// Instance
func (h *Handlers) InstCreate(ctx echo.Context) error {
	return h.InstHandler.Create(echoc.NewContext(ctx))
}

func (h *Handlers) InstPut(ctx echo.Context) error {
	return h.InstHandler.Put(echoc.NewContext(ctx))
}

func (h *Handlers) InstGet(ctx echo.Context) error {
	return h.InstHandler.Get(echoc.NewContext(ctx))
}

func (h *Handlers) InstListSpaces(ctx echo.Context) error {
	return h.InstHandler.ListSpaces(echoc.NewContext(ctx))
}

// Space
func (h *Handlers) SpaceCreate(ctx echo.Context) error {
	return h.SpaceHandler.Create(echoc.NewContext(ctx))
}

func (h *Handlers) SpacePut(ctx echo.Context) error {
	return h.SpaceHandler.Put(echoc.NewContext(ctx))
}

func (h *Handlers) SpaceGet(ctx echo.Context) error {
	return h.SpaceHandler.Get(echoc.NewContext(ctx))
}

// Ente
func (h *Handlers) EnteCreate(ctx echo.Context) error {
	return h.EnteHandler.Create(echoc.NewContext(ctx))
}

func (h *Handlers) EntePut(ctx echo.Context) error {
	return h.EnteHandler.Put(echoc.NewContext(ctx))
}

func (h *Handlers) EnteGet(ctx echo.Context) error {
	return h.EnteHandler.Get(echoc.NewContext(ctx))
}
