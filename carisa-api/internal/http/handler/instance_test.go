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
	"testing"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/api/internal/instance"

	"github.com/carisa/api/internal/http"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestInstance_Create(t *testing.T) {
	tests := []struct {
		name     string
		error    bool
		body     string
		status   uint16
		response string
	}{
		{
			name:     "Body wrong. Bad request",
			error:    true,
			body:     "{df",
			status:   nethttp.StatusBadRequest,
			response: "",
		},
	}

	e := echo.New()
	handler, mng := NewHandler(t)
	defer mng.Close()

	for _, tt := range tests {
		rec, ctx := http.MockHttp(e, "/api/instances", tt.body)
		err := handler.Create(ctx)
		if tt.error {
			assert.Error(t, err, tt.name)
		} else {
			if assert.NoError(t, err, tt.name) {
				assert.Equal(t, tt.status, rec.Code, tt.name)
				assert.Equal(t, tt.response, rec.Body.String(), tt.name)
			}
		}
	}
}

func NewHandler(t *testing.T) (Instance, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt := mock.NewContainerFake()
	srv := instance.NewService(cnt, mng.Store())
	return NewInstanceHandl(srv, cnt), mng
}
