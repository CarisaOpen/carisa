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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/api/mock"

	"github.com/stretchr/testify/assert"

	"github.com/rs/xid"
)

func TestConverter_ParamID(t *testing.T) {
	h := mock.HTTP()
	id := xid.New()
	defer h.Close(nil)

	params := map[string]string{
		"id": id.String(),
	}

	_, ctx := h.NewHTTP(http.MethodGet, "/api/:id", "", params, nil)
	idr, err := ParamID(ctx)

	if assert.NoError(t, err) {
		assert.Equal(t, idr, id)
	}
}

func TestConverter_ParamID_MissingParamError(t *testing.T) {
	h := mock.HTTP()
	defer h.Close(nil)

	_, ctx := h.NewHTTP(http.MethodGet, "/api", "", map[string]string{}, nil)
	_, err := ParamID(ctx)

	assert.Error(t, err)
}

func TestConverter_ParamID_ConvertParamError(t *testing.T) {
	h := mock.HTTP()
	defer h.Close(nil)
	params := map[string]string{
		"id": "123)",
	}

	_, ctx := h.NewHTTP(http.MethodGet, "/api/id", "", params, nil)
	_, err := ParamID(ctx)

	assert.Error(t, err)
}

func TestFilterLink(t *testing.T) {
	tests := []struct {
		qparam map[string]string
		id     string
		pname  string
		name   string
		top    int
		error  bool
		ranges bool
		exclID bool
	}{
		{
			name:   "List start name.",
			error:  false,
			id:     xid.New().String(),
			pname:  "sname",
			top:    1,
			ranges: false,
			qparam: map[string]string{
				"sname": "sname",
				"top":   "1",
			},
		},
		{
			name:   "Range.",
			error:  false,
			id:     xid.New().String(),
			pname:  "gtname",
			top:    20,
			ranges: true,
			qparam: map[string]string{
				"gtname": "gtname",
			},
		},
		{
			name:   "Range. Exclude ID.",
			error:  false,
			id:     xid.NilID().String(),
			pname:  "gtname",
			top:    20,
			ranges: true,
			qparam: map[string]string{
				"gtname": "gtname",
			},
		},
		{
			name:   "ID, incorrect format.",
			error:  true,
			id:     "123",
			pname:  "",
			top:    0,
			ranges: false,
			qparam: map[string]string{
				"top": "top",
			},
		},
		{
			name:   "Parameters missing.",
			error:  true,
			id:     xid.NilID().String(),
			pname:  "",
			top:    0,
			ranges: false,
			qparam: map[string]string{},
		},
		{
			name:   "Top, incorrect format.",
			error:  true,
			id:     xid.NilID().String(),
			pname:  "",
			top:    0,
			ranges: false,
			qparam: map[string]string{
				"top": "top",
			},
		},
		{
			name:   "Top: 0. Between 1 and 100",
			error:  true,
			id:     xid.NilID().String(),
			pname:  "",
			top:    0,
			ranges: false,
			qparam: map[string]string{
				"top": "0",
			},
		},
		{
			name:   "Top: 101. Between 1 and 100",
			error:  true,
			id:     xid.NilID().String(),
			pname:  "",
			top:    0,
			ranges: false,
			qparam: map[string]string{
				"top": "101",
			},
		},
	}

	h := mock.HTTP()
	defer h.Close(nil)

	for _, tt := range tests {
		_, ctx := h.NewHTTP(http.MethodGet, "/api/:id", "", map[string]string{"id": tt.id}, tt.qparam)
		id, name, top, ranges, err := FilterLink(ctx, tt.exclID)
		if tt.error {
			assert.Error(t, err, strings.Concat(tt.name, "Error"))
		} else {
			assert.Equal(t, id.String(), tt.id, strings.Concat(tt.name, "Id"))
		}
		assert.Equal(t, name, tt.pname, strings.Concat(tt.name, "Name"))
		assert.Equal(t, top, tt.top, strings.Concat(tt.name, "Top"))
		assert.Equal(t, ranges, tt.ranges, strings.Concat(tt.name, "Range"))
	}
}
