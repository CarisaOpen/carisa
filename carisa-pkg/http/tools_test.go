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
	httpc "net/http"
	"testing"

	"github.com/carisa/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestHttpTools_CreateStatus(t *testing.T) {
	tests := []struct {
		found  bool
		status int
	}{
		{
			found:  true,
			status: httpc.StatusCreated,
		},
		{
			found:  false,
			status: httpc.StatusFound,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.status, CreateStatus(tt.found))
	}
}

func TestHttpTools_PutStatus(t *testing.T) {
	tests := []struct {
		updated bool
		status  int
	}{
		{
			updated: true,
			status:  httpc.StatusOK,
		},
		{
			updated: false,
			status:  httpc.StatusCreated,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.status, PutStatus(tt.updated))
	}
}

func TestHttpTools_GetStatus(t *testing.T) {
	tests := []struct {
		found  bool
		status int
	}{
		{
			found:  true,
			status: httpc.StatusOK,
		},
		{
			found:  false,
			status: httpc.StatusNotFound,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.status, GetStatus(tt.found))
	}
}

func TestHttpTools_Close(t *testing.T) {
	e := echo.New()
	Close(newLogger(zap.DebugLevel), e)
}

func newLogger(level zapcore.Level) logging.Logger {
	core, _ := observer.New(level)
	return logging.NewZapWrap(zap.New(core), logging.DebugLevel, "")
}
