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
	"errors"
	"net/http"
	"testing"

	"github.com/carisa/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContext_Param(t *testing.T) {
	e := echo.New()
	defer closeEcho(e)

	params := map[string]string{
		"param1": "value1",
	}
	_, ctx := MockHTTP(e, http.MethodGet, "/api", "", params)
	ctxw := NewContext(ctx)

	value, _ := ctxw.Param("param1")
	assert.Equal(t, "value1", value)
	_, err := ctxw.Param("param2")
	assert.Error(t, err)
}

func TestContext_Bind(t *testing.T) {
	s := struct {
		P int `json:"p,omitempty"`
	}{
		P: 1,
	}

	e := echo.New()
	defer closeEcho(e)

	_, ctx := MockHTTP(e, http.MethodPost, "/api", `{"p":1}`, nil)
	ctxw := NewContext(ctx)

	if assert.NoError(t, ctxw.Bind(&s)) {
		assert.Equal(t, s.P, 1)
	}
}

func TestContext_JSON(t *testing.T) {
	s := struct {
		P int `json:"p,omitempty"`
	}{
		P: 1,
	}

	e := echo.New()
	defer closeEcho(e)
	body := `{"p":1}`

	rec, ctx := MockHTTP(e, http.MethodPost, "/api", body, nil)
	ctxw := NewContext(ctx)

	if assert.NoError(t, ctxw.JSON(200, s)) {
		assert.Contains(t, rec.Body.String(), body)
	}
}

func TestContext_HTTPErrorLog(t *testing.T) {
	recorded, l := newLogger(zapcore.ErrorLevel)
	ctxw := NewContext(nil)

	err := ctxw.HTTPErrorLog(
		500,
		"error",
		errors.New("parent error"),
		l,
		"loc",
		logging.String("k", "v"))

	assert.Equal(t, "code=500, message=error. k: v", err.Error(), "error")

	for _, logs := range recorded.All() {
		assert.Equal(t, "error. k: v: parent error", logs.Message, "log")
	}
}

func TestContext_HTTPError(t *testing.T) {
	ctxw := NewContext(nil)
	err := ctxw.HTTPError(500, "error")
	assert.Equal(t, "code=500, message=[error]", err.Error(), "error")
}

func TestValid_ValidNoEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		message string
	}{
		{
			name:    "property",
			value:   "value",
			message: "",
		},
		{
			name:    "property",
			value:   "",
			message: "code=400, message=[the property: 'property' can not be empty]",
		},
	}

	ctxw := NewContext(nil)

	for _, tt := range tests {
		r := ctxw.NoEmpty(tt.name, tt.value)
		if len(tt.message) == 0 {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, tt.message, r.Error())
		}
	}
}

func newLogger(level zapcore.Level) (*observer.ObservedLogs, logging.Logger) {
	core, obs := observer.New(level)
	return obs, logging.NewZapWrap(zap.New(core), logging.DebugLevel, "")
}

func closeEcho(e *echo.Echo) {
	_ = e.Close()
}
