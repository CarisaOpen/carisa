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
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/carisa/pkg/logging"

	"github.com/stretchr/testify/assert"
)

func TestCreateStatus(t *testing.T) {
	tests := []struct {
		found  bool
		status int
	}{
		{
			found:  true,
			status: http.StatusCreated,
		},
		{
			found:  false,
			status: http.StatusFound,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.status, CreateStatus(tt.found))
	}
}

func TestNewHTTPErrorLog(t *testing.T) {
	log, err := logging.NewZapWrapDev()
	if err != nil {
		t.Error(err)
	}

	errLog := NewHTTPErrorLog(
		http.StatusBadRequest,
		"error",
		errors.New("error parent"),
		log,
		"loc",
		logging.String("key", "value"))

	assert.Equal(t, http.StatusBadRequest, errLog.Code, "Error status")
	assert.Equal(t, "error. key: value", errLog.Message, "Message")
}

func TestClose(t *testing.T) {
	e := echo.New()
	Close(e)
}
