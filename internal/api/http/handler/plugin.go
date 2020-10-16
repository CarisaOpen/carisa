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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/api/plugin"

	"github.com/carisa/internal/api/http/convert"

	"github.com/carisa/pkg/http"

	"github.com/carisa/internal/api/runtime"
	httpc "github.com/carisa/pkg/http"
)

const locPlugin = "http.plugin"

// Plugin hands the http request of the plugin
type Plugin struct {
	srv plugin.Service
	cnt *runtime.Container
}

// NewPluginHandle creates handler
func NewPluginHandle(srv plugin.Service, cnt *runtime.Container) Plugin {
	return Plugin{
		srv: srv,
		cnt: cnt,
	}
}

// Create creates the plugin
func (p *Plugin) Create(c httpc.Context, category plugin.Category) error {
	proto := plugin.Prototype{}
	if err := bind(c, locPlugin, p.cnt.Log, &proto); err != nil {
		return err
	}
	proto.Category = category

	created, err := p.srv.Create(&proto)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to create the plugin prototype")
	}

	return c.JSON(http.CreateStatus(created), proto)
}

// Put creates or update the plugin
func (p *Plugin) Put(c httpc.Context, category plugin.Category) error {
	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	proto := plugin.Prototype{}
	if err := bind(c, locPlugin, p.cnt.Log, &proto); err != nil {
		return err
	}
	proto.Category = category

	proto.ID = id
	updated, err := p.srv.Put(&proto)
	if err != nil {
		return c.HTTPError(
			nethttp.StatusInternalServerError,
			"it was impossible to create or update the plugin prototype")
	}

	return c.JSON(http.PutStatus(updated), proto)
}

// Get gets the plugin by ID
func (p *Plugin) Get(c httpc.Context) error {
	var proto plugin.Prototype

	id, err := convert.ParamID(c)
	if err != nil {
		return err
	}

	found, err := p.srv.Get(id, &proto)
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, "it was impossible to get the plugin prototype")
	}

	return c.JSON(http.GetStatus(found), proto)
}

// ListProps list properties by ente ID and return top properties.
// If sname query param is not empty, is filtered by properties which name starts by name parameter
// If gtname query param is not empty, is filtered by properties which name is greater than name parameter
func (p *Plugin) ListPlugins(c httpc.Context, cat plugin.Category) error {
	_, name, top, ranges, err := convert.FilterLink(c, true)
	if err != nil {
		return err
	}

	props, err := p.srv.ListPlugins(cat, name, ranges, top)
	if err != nil {
		return c.HTTPError(
			nethttp.StatusInternalServerError,
			strings.Concat("it was impossible to list the plugins (", string(cat), ")"))
	}

	return c.JSON(nethttp.StatusOK, props)
}
