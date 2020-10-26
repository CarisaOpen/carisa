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
	"github.com/carisa/internal/api/plugin"
	echoc "github.com/carisa/pkg/http/echo"
	"github.com/labstack/echo/v4"
)

// Handlers is a handlers container
type Handlers struct {
	InstHandler     Instance
	SpaceHandler    Space
	EnteHandler     Ente
	CategoryHandler Category
	PluginHandler   Plugin
	ObjectHandler   Object
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

func (h *Handlers) SpcListEntes(ctx echo.Context) error {
	return h.SpaceHandler.ListEntes(echoc.NewContext(ctx))
}

func (h *Handlers) SpcListCategories(ctx echo.Context) error {
	return h.SpaceHandler.ListCategories(echoc.NewContext(ctx))
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

func (h *Handlers) EnteListProps(ctx echo.Context) error {
	return h.EnteHandler.ListProps(echoc.NewContext(ctx))
}

func (h *Handlers) EnteLinkToCat(ctx echo.Context) error {
	return h.EnteHandler.LinkToCat(echoc.NewContext(ctx))
}

func (h *Handlers) EnteCreateProp(ctx echo.Context) error {
	return h.EnteHandler.CreateProp(echoc.NewContext(ctx))
}

func (h *Handlers) EntePutProp(ctx echo.Context) error {
	return h.EnteHandler.PutProp(echoc.NewContext(ctx))
}

func (h *Handlers) EnteGetProp(ctx echo.Context) error {
	return h.EnteHandler.GetProp(echoc.NewContext(ctx))
}

// Category
func (h *Handlers) CatCreate(ctx echo.Context) error {
	return h.CategoryHandler.Create(echoc.NewContext(ctx))
}

func (h *Handlers) CatPut(ctx echo.Context) error {
	return h.CategoryHandler.Put(echoc.NewContext(ctx))
}

func (h *Handlers) CatGet(ctx echo.Context) error {
	return h.CategoryHandler.Get(echoc.NewContext(ctx))
}

func (h *Handlers) CatListCategories(ctx echo.Context) error {
	return h.CategoryHandler.ListCategories(echoc.NewContext(ctx))
}

func (h *Handlers) CatListProps(ctx echo.Context) error {
	return h.CategoryHandler.ListProps(echoc.NewContext(ctx))
}

func (h *Handlers) CatGetProp(ctx echo.Context) error {
	return h.CategoryHandler.GetProp(echoc.NewContext(ctx))
}

func (h *Handlers) CatCreateProp(ctx echo.Context) error {
	return h.CategoryHandler.CreateProp(echoc.NewContext(ctx))
}

func (h *Handlers) CatPutProp(ctx echo.Context) error {
	return h.CategoryHandler.PutProp(echoc.NewContext(ctx))
}

func (h *Handlers) CatPropLinkTo(ctx echo.Context) error {
	return h.CategoryHandler.LinkToProp(echoc.NewContext(ctx))
}

// Plugin query prototype
func (h *Handlers) PluginQryCreate(ctx echo.Context) error {
	return h.PluginHandler.Create(echoc.NewContext(ctx), plugin.Query)
}

func (h *Handlers) PluginQryPut(ctx echo.Context) error {
	return h.PluginHandler.Put(echoc.NewContext(ctx), plugin.Query)
}

func (h *Handlers) PluginQryGet(ctx echo.Context) error {
	return h.PluginHandler.Get(echoc.NewContext(ctx))
}

func (h *Handlers) PluginQryListPlugins(ctx echo.Context) error {
	return h.PluginHandler.ListPlugins(echoc.NewContext(ctx), plugin.Query)
}

// Query object instance
func (h *Handlers) InstQryCreate(ctx echo.Context) error {
	return h.ObjectHandler.Create(echoc.NewContext(ctx), plugin.Query)
}

func (h *Handlers) InstQryPut(ctx echo.Context) error {
	return h.ObjectHandler.Put(echoc.NewContext(ctx), plugin.Query)
}

func (h *Handlers) InstQryGet(ctx echo.Context) error {
	return h.ObjectHandler.Get(echoc.NewContext(ctx))
}

func (h *Handlers) InstQryListQueries(ctx echo.Context) error {
	return h.ObjectHandler.ListInstances(echoc.NewContext(ctx), plugin.Query)
}
