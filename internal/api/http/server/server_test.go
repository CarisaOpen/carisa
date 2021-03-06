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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/carisa/internal/api/http/handler"

	"github.com/labstack/echo/v4"
)

func TestServer_Middleware(t *testing.T) {
	e := echo.New()
	Middleware(e)
}

func TestServer_Router(t *testing.T) {
	e := echo.New()
	h := handler.Handlers{}

	Router(e, h)

	assert.Equal(t, 37, len(e.Routes()))
}
