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

package plugin

import (
	"testing"

	"github.com/carisa/internal/api/samples"

	"github.com/carisa/internal/api/test"

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.plugin_test

func TestPluginService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	proto := proto()

	ok, err := srv.Create(proto)

	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
		checkProto(t, srv, "Checking relations", *proto)
	}
}

func TestPluginService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	tests := []struct {
		name    string
		updated bool
		proto   *Prototype
	}{
		{
			name:    "Creating Plugin",
			updated: false,
			proto: &Prototype{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				Category: Query,
			},
		},
		{
			name:    "Updating Plugin",
			updated: true,
			proto: &Prototype{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				Category: Query,
			},
		},
	}

	for _, tt := range tests {
		updated, err := srv.Put(tt.proto)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Plugin updated"))
			checkProto(t, srv, tt.name, *tt.proto)
		}
	}
}

func checkProto(t *testing.T, srv Service, name string, proto Prototype) {
	var pr Prototype
	_, err := srv.Get(proto.ID, &pr)
	if assert.NoError(t, err) {
		assert.Equal(t, proto, pr, "Getting plugin")
	}
	test.CheckRelations(t, srv.crud.Store(), name, &proto)
}

func TestPluginService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	proto := proto()

	_, err := srv.Create(proto)
	if assert.NoError(t, err) {
		var get Prototype
		ok, err := srv.Get(proto.ID, &get)
		if assert.NoError(t, err) {
			assert.True(t, ok, "Get ok")
			assert.Equal(t, proto, &get, "Plugin returned")
		}
	}
}

func TestPluginService_ListPlugins(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	proto := New()
	proto.Name = "nameproto"
	link := proto.Link()

	_, err := s.crud.Create("", s.cnt.StoreWithTimeout, link)

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := s.ListPlugins(Query, "nameproto", tt.Ranges, 1)
			if assert.NoError(t, err, tt.Name) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
			}
		}
	}
}

func TestPluginService_Exists(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	proto := proto()

	_, err := srv.Create(proto)
	if assert.NoError(t, err) {
		found, err := srv.Exists(proto.ID)
		if assert.NoError(t, err) {
			assert.True(t, found, "Exists")
		}
	}
}

func TestPluginService_ExistsWithError(t *testing.T) {
	srv, crud := newServiceMocked()
	crud.Store().(*storage.ErrMockCRUD).Activate("Exists")

	_, err := srv.Exists(xid.New())
	assert.Error(t, err)
}

func proto() *Prototype {
	proto := New()
	proto.Name = "name"
	proto.Desc = "desc"
	return &proto
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crudOper.Store())
	return NewService(cnt, ext, crudOper), mng
}

func newServiceMocked() (Service, *storage.ErrMockCRUDOper) {
	cnt := mock.NewContainerFake()
	crud := storage.NewErrMockCRUDOper()
	ext := service.NewExt(cnt, crud.Store())
	return NewService(cnt, ext, crud), crud
}
