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

package valid

import (
	"net/http"

	"github.com/carisa/pkg/strings"
	"github.com/labstack/echo/v4"
)

func NoEmpty(name string, value string) *echo.HTTPError {
	if len(value) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, strings.Concat("the property: '", name, "' can not be empty"))
	}
	return nil
}
