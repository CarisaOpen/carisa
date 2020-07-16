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

	"github.com/carisa/api/internal/http/handler"
	"github.com/carisa/api/internal/runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Middleware configure security and behaviour of http
func Middleware(e *echo.Echo) {
	e.Use(middleware.Recover())
}

// Router defines all http route for API
func Router(e *echo.Echo, hands handler.Handlers) {
	// InstCreate
	e.POST("/api/instances", hands.InstCreate)
}

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
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error(err.Error())
	}
}
