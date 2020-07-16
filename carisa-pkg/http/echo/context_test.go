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
	"testing"

	"github.com/carisa/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContext_Bind(t *testing.T) {
	s := struct {
		P int `json:"p,omitempty"`
	}{
		P: 1,
	}

	e := echo.New()
	defer closeEcho(e)

	_, ctx := MockHTTPPost(e, "/api/instances", `{"p":1}`)
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

	rec, ctx := MockHTTPPost(e, "/api/instances", body)
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

func newLogger(level zapcore.Level) (*observer.ObservedLogs, logging.Logger) {
	core, obs := observer.New(level)
	return obs, logging.NewZapWrap(zap.New(core), logging.DebugLevel, "")
}

func closeEcho(e *echo.Echo) {
	_ = e.Close()
}
