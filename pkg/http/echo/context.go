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
	"strconv"

	"github.com/carisa/pkg/http"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/strings"
	"github.com/labstack/echo/v4"

	nethttp "net/http"
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

// Param implements Context.Param
func (c *context) Param(name string) string {
	return c.ctx.Param(name)
}

// QueryParam implements Context.QueryParam
func (c *context) QueryParam(name string) string {
	return c.ctx.QueryParam(name)
}

// Bind implements Context.Bind
func (c *context) Bind(i interface{}) error {
	return c.ctx.Bind(i)
}

// JSON implements Context.JSON
func (c *context) JSON(code int, i interface{}) error {
	return c.ctx.JSON(code, i)
}

// HTTPErrorLog implements Context.HTTPErrorLog
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

// HTTPError implements Context.HTTPError
func (c *context) HTTPError(code int, message ...interface{}) error {
	return echo.NewHTTPError(code, message)
}

// NoEmpty implements Context.NoEmpty
func (c *context) NoEmpty(name string, value string) error {
	if len(value) == 0 {
		return c.HTTPError(nethttp.StatusBadRequest, strings.Concat("the property: '", name, "' can not be empty"))
	}
	return nil
}

// MaxLen implements Context.MaxLen
func (c *context) MaxLen(name string, value string, length int) error {
	if len(value) > length {
		return c.HTTPError(
			nethttp.StatusBadRequest,
			strings.Concat("the property: '", name, "' can not be more than ", strconv.Itoa(length)))
	}
	return nil
}
