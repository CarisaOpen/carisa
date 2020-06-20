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

package http

import (
	nethttp "net/http"

	"github.com/carisa/pkg/logging"
	"github.com/labstack/echo/v4"
)

// NewHTTPErrorLog creates http error and sending a log error
func NewHTTPErrorLog(
	status int,
	msg string, err error,
	logger logging.Logger,
	loc string, fields ...logging.Field) *echo.HTTPError {
	_ = logger.ErrWrap(err, msg, loc, fields...)
	return echo.NewHTTPError(status, logging.Compose(msg, fields...))
}

// status return status for create handler
func CreateStatus(found bool) int {
	status := nethttp.StatusCreated
	if found {
		status = nethttp.StatusFound
	}
	return status
}
