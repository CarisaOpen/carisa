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

package validator

import (
	"strings"
	"testing"

	"github.com/rs/xid"

	"github.com/carisa/api/internal/mock"

	"github.com/carisa/api/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestValid_ValidDescriptor(t *testing.T) {
	tests := []struct {
		desc    entity.Descriptor
		message string
	}{
		{
			desc: entity.Descriptor{
				Name: "",
				Desc: "desc",
			},
			message: "code=400, message=[the property: 'name' can not be empty]",
		},
		{
			desc: entity.Descriptor{
				Name: "name",
				Desc: "",
			},
			message: "code=400, message=[the property: 'description' can not be empty]",
		},
		{
			desc: entity.Descriptor{
				Name: strings.Repeat("n", 51),
				Desc: "desc",
			},
			message: "code=400, message=[the property: 'name' can not be more than 50]",
		},
		{
			desc: entity.Descriptor{
				Name: "name",
				Desc: strings.Repeat("d", 501),
			},
			message: "code=400, message=[the property: 'description' can not be more than 500]",
		},
		{
			desc: entity.Descriptor{
				Name: "name",
				Desc: "desc",
			},
			message: "",
		},
	}

	ctx := mock.NewContextFake()

	for _, tt := range tests {
		r := Descriptor(ctx, tt.desc)
		if len(tt.message) == 0 {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, tt.message, r.Error())
		}
	}
}

func TestValid_ValidID(t *testing.T) {
	tests := []struct {
		id      xid.ID
		message string
	}{
		{
			message: "code=400, message=[the property: 'ID' can not be empty]",
		},
		{
			id:      xid.New(),
			message: "",
		},
	}

	ctx := mock.NewContextFake()

	for _, tt := range tests {
		r := ID(ctx, tt.id)
		if len(tt.message) == 0 {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, tt.message, r.Error())
		}
	}
}
