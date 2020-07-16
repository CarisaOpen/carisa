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

package echo

import (
	"github.com/carisa/pkg/http"
	"github.com/carisa/pkg/logging"
	"github.com/labstack/echo/v4"
)

// NewContext creates the echo context adapter
func NewContext(ctx echo.Context) http.Context {
	return &context{
		ctx: ctx,
	}
}

// Context is a adapter for echo
type context struct {
	ctx echo.Context
}

// Bind implements http.interface.Context.Bind
func (c *context) Bind(i interface{}) error {
	return c.ctx.Bind(i)
}

// JSON implements http.interface.Context.JSON
func (c *context) JSON(code int, i interface{}) error {
	return c.ctx.JSON(code, i)
}

// HTTPErrorLog implements http.interface.Context.HTTPErrorLog
func (c *context) HTTPErrorLog(
	status int,
	msg string, err error,
	logger logging.Logger,
	loc string, fields ...logging.Field) error {
	switch len(fields) {
	case 0:
		_ = logger.ErrWrap(err, msg, loc)
	case 1:
		_ = logger.ErrWrap1(err, msg, loc, fields[0])
	case 2:
		_ = logger.ErrWrap2(err, msg, loc, fields[0], fields[1])
	case 3:
		_ = logger.ErrWrap3(err, msg, loc, fields[0], fields[1], fields[2])
	}
	return echo.NewHTTPError(status, logging.Compose(msg, fields...))
}

// HTTPError implements http.interface.Context.HTTPError
func (c *context) HTTPError(code int, message ...interface{}) error {
	return echo.NewHTTPError(code, message)
}
