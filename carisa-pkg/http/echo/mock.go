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
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

func MockHTTP(
	e *echo.Echo, method string,
	url string,
	body string,
	params map[string]string) (rec *httptest.ResponseRecorder, c echo.Context) {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	setParams(params, c)
	return
}

func setParams(params map[string]string, c echo.Context) {
	if params == nil {
		return
	}

	l := len(params)
	names := make([]string, 0, l)
	values := make([]string, 0, l)
	for k, v := range params {
		names = append(names, k)
		values = append(values, v)
	}

	c.SetParamNames(names...)
	c.SetParamValues(values...)
}
