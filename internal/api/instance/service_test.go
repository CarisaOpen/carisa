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

package instance

import (
	"testing"

	srv "github.com/carisa/internal/api/service"

	"github.com/carisa/internal/api/samples"
	spacesmpl "github.com/carisa/internal/api/space/samples"
	"github.com/rs/xid"

	"github.com/carisa/pkg/storage"

	"github.com/carisa/internal/api/mock"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.instance_test

func TestInstanceService_Create(t *testing.T) {
	i := instance()
	s, mng := newServiceFaked(t)
	defer mng.Close()

	ok, err := s.Create(&i)

	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
		checkInstance(t, s, i)
	}
}

func TestInstanceService_Put(t *testing.T) {
	i := instance()
	s, mng := newServiceFaked(t)
	defer mng.Close()

	i.AutoID()
	ok, err := s.Put(&i)

	if assert.NoError(t, err) {
		assert.False(t, ok, "Created")
		checkInstance(t, s, i)
	}
}

func checkInstance(t *testing.T, s Service, i Instance) {
	var ir Instance
	_, err := s.Get(i.ID, &ir)
	if assert.NoError(t, err) {
		assert.Equal(t, i, ir, "Getting instance")
	}
}

func TestInstanceService_Get(t *testing.T) {
	i := instance()
	s, mng := newServiceFaked(t)
	defer mng.Close()

	_, err := s.Create(&i)
	if assert.NoError(t, err) {
		var geti Instance
		ok, err := s.Get(i.ID, &geti)
		if assert.NoError(t, err) {
			assert.True(t, ok, "Get ok")
			assert.Equal(t, i, geti, "Instance returned")
		}
	}
}

func TestInstanceService_ListSpaces(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	link, _, err := spacesmpl.CreateLink(mng, id)

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := s.ListSpaces(id, "name", tt.Ranges, 1)
			if assert.NoError(t, err) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
			}
		}
	}
}

func instance() Instance {
	inst := New()
	inst.Name = "name"
	inst.Desc = "desc"
	return inst
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	ext := srv.NewExt(cnt, crudOper.Store())
	return NewService(cnt, ext, crudOper), mng
}
