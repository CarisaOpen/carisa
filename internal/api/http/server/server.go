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

package server

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/carisa/internal/api/http/handler"
	"github.com/carisa/internal/api/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start starts the graceful http server
func Start(e *echo.Echo, cnf runtime.Config) {
	// Start server
	go func() {
		if err := e.Start(cnf.Server.Address()); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error(err.Error())
	}
}

// Middleware configure security and behaviour of http
func Middleware(e *echo.Echo) {
	e.Use(middleware.Recover())
}

// Router defines all http route for API
func Router(e *echo.Echo, h handler.Handlers) {
	// Instance
	e.POST("/api/instances", h.InstCreate)
	e.PUT("/api/instances/:id", h.InstPut)
	e.GET("/api/instances/:id", h.InstGet)
	e.GET("/api/instances/:id/spaces", h.InstListSpaces)

	// Space
	e.POST("/api/spaces", h.SpaceCreate)
	e.PUT("/api/spaces/:id", h.SpacePut)
	e.GET("/api/spaces/:id", h.SpaceGet)
	e.GET("/api/spaces/:id/entes", h.SpcListEntes)
	e.GET("/api/spaces/:id/categories", h.SpcListCategories)

	// Ente
	e.POST("/api/entes", h.EnteCreate)
	e.PUT("/api/entes/:id", h.EntePut)
	e.GET("/api/entes/:id", h.EnteGet)
	e.GET("/api/entes/:id/properties", h.EnteListProps)
	e.POST("/api/entes/:id/queries", h.EnteQryCreate)
	e.PUT("/api/entes/:enteid/queries/:id", h.EnteQryPut)
	e.GET("/api/entes/:id/queries", h.EnteListQueries)
	e.PUT("/api/entes/:enteid/linktocategories/:categoryid", h.EnteLinkToCat)
	e.POST("/api/entesproperties", h.EnteCreateProp)
	e.PUT("/api/entesproperties/:id", h.EntePutProp)
	e.GET("/api/entesproperties/:id", h.EnteGetProp)

	// Category
	e.POST("/api/categories", h.CatCreate)
	e.PUT("/api/categories/:id", h.CatPut)
	e.GET("/api/categories/:id", h.CatGet)
	e.GET("/api/categories/:id/child", h.CatListCategories)
	e.GET("/api/categories/:id/properties", h.CatListProps)
	e.POST("/api/categories/:id/queries", h.CatQryCreate)
	e.PUT("/api/categories/:categoryid/queries/:id", h.CatQryPut)
	e.GET("/api/categories/:id/queries", h.CatListQueries)
	e.POST("/api/categoriesproperties", h.CatCreateProp)
	e.PUT("/api/categoriesproperties/:id", h.CatPutProp)
	e.GET("/api/categoriesproperties/:id", h.CatGetProp)
	e.PUT("/api/categoriesproperties/:catpropid/linkto/:propid", h.CatPropLinkTo)

	// Query plugin Prototype
	e.POST("/api/plugins/queries", h.PluginQryCreate)
	e.PUT("/api/plugins/queries/:id", h.PluginQryPut)
	e.GET("/api/plugins/queries/:id", h.PluginQryGet)
	e.GET("/api/plugins/queries", h.PluginQryListPlugins)

	// Query object Instance
	e.GET("/api/queries/:id", h.InstQryGet)
}
