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

package convert

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	echoc "github.com/carisa/pkg/http/echo"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

func TestConverter_ParamID(t *testing.T) {
	e := echo.New()
	id := xid.New()
	params := map[string]string{
		"id": id.String(),
	}

	_, ctx := echoc.MockHTTP(e, http.MethodGet, "/api/:id", "", params)
	idr, err := ParamID(echoc.NewContext(ctx))

	if assert.NoError(t, err) {
		assert.Equal(t, idr, id)
	}
}

func TestConverter_ParamID_MissingParamError(t *testing.T) {
	e := echo.New()

	_, ctx := echoc.MockHTTP(e, http.MethodGet, "/api", "", map[string]string{})
	_, err := ParamID(echoc.NewContext(ctx))

	assert.Error(t, err)
}

func TestConverter_ParamID_ConvertParamError(t *testing.T) {
	e := echo.New()
	params := map[string]string{
		"id": "123)",
	}

	_, ctx := echoc.MockHTTP(e, http.MethodGet, "/api/id", "", params)
	_, err := ParamID(echoc.NewContext(ctx))

	assert.Error(t, err)
}
