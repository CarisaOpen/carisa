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

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/http/validator"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/strings"

	httpc "github.com/carisa/pkg/http"
)

// errCRUDSrv checks the service errors
func errCRUDSrv(c httpc.Context, err error, msg string, msgNotFound string, found bool) error {
	if err != nil {
		return c.HTTPError(nethttp.StatusInternalServerError, msg)
	}
	if !found {
		return c.HTTPError(nethttp.StatusNotFound, msgNotFound)
	}
	return nil
}

// bind binds entity from http body and doing validation
func bind(c httpc.Context, loc string, log logging.Logger, e entity.Domain) error {
	if err := c.Bind(e); err != nil {
		return c.HTTPErrorLog(
			nethttp.StatusBadRequest,
			strings.Concat("cannot recover ", e.ToString()),
			err,
			log,
			loc)
	}

	if httpErr := validator.Descriptor(c, e.Nominative()); httpErr != nil {
		return httpErr
	}
	return nil
}
