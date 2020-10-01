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

	"github.com/stretchr/testify/assert"
)

func TestContext_Param(t *testing.T) {
	h := HTTPMock()
	defer h.Close(nil)

	params := map[string]string{
		"param1": "value1",
	}
	_, ctx := h.NewHTTP(http.MethodGet, "/api", "", params, nil)

	assert.Equal(t, "value1", ctx.Param("param1"))
}

func TestContext_QueryParam(t *testing.T) {
	h := HTTPMock()
	defer h.Close(nil)

	params := map[string]string{
		"param1": "value1",
	}
	_, ctx := h.NewHTTP(http.MethodGet, "/api", "", nil, params)

	assert.Equal(t, "value1", ctx.QueryParam("param1"))
}

func TestContext_Bind(t *testing.T) {
	s := struct {
		P int `json:"p,omitempty"`
	}{
		P: 1,
	}

	h := HTTPMock()
	defer h.Close(nil)

	_, ctx := h.NewHTTP(http.MethodPost, "/api", `{"p":1}`, nil, nil)

	if assert.NoError(t, ctx.Bind(&s)) {
		assert.Equal(t, s.P, 1)
	}
}

func TestContext_JSON(t *testing.T) {
	s := struct {
		P int `json:"p,omitempty"`
	}{
		P: 1,
	}

	h := HTTPMock()
	defer h.Close(nil)
	body := `{"p":1}`

	rec, ctx := h.NewHTTP(http.MethodPost, "/api", body, nil, nil)

	if assert.NoError(t, ctx.JSON(200, s)) {
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

func TestValid_ValidMaxLen(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		len     int
		message string
	}{
		{
			name:    "property",
			value:   "value",
			len:     10,
			message: "",
		},
		{
			name:    "property",
			value:   "value",
			len:     3,
			message: "code=400, message=[the property: 'property' can not be more than 3]",
		},
	}

	ctxw := NewContext(nil)

	for _, tt := range tests {
		r := ctxw.MaxLen(tt.name, tt.value, tt.len)
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
