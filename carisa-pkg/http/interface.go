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

import "github.com/carisa/pkg/logging"

// Context is a adapter for http
type Context interface {
	// Param return path param
	Param(name string) (string, error)

	// Bind binds the request body into provided type `i`. The default binder
	// does it based on Content-Type header.
	Bind(i interface{}) error

	// JSON sends a JSON response with status code.
	JSON(code int, i interface{}) error

	// HTTPErrorLog creates http error and sending a log error
	HTTPErrorLog(status int, msg string, err error, logger logging.Logger, loc string, fields ...logging.Field) error

	// HTTPError return http error
	HTTPError(code int, message ...interface{}) error

	// NoEmpty validates that the name is not empty
	NoEmpty(name string, value string) error
}
