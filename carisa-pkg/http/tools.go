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
	"net/http"

	"github.com/carisa/pkg/logging"
	"github.com/labstack/echo/v4"
)

// CreateStatus return status for create handler
func CreateStatus(created bool) int {
	status := http.StatusFound
	if created {
		status = http.StatusCreated
	}
	return status
}

// PutStatus return status for put handler
func PutStatus(updated bool) int {
	status := http.StatusCreated
	if updated {
		status = http.StatusOK
	}
	return status
}

// Close close echo connection
func Close(log logging.Logger, e *echo.Echo) {
	if err := e.Close(); err != nil {
		log.PanicE(err, "http.close")
	}
}
