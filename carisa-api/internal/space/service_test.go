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

package space

import (
	"testing"

	entesmpl "github.com/carisa/api/internal/ente/samples"
	"github.com/carisa/api/internal/samples"
	srv "github.com/carisa/api/internal/service"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/entity"
	"github.com/rs/xid"

	catsmpl "github.com/carisa/api/internal/category/samples"
	instsmpl "github.com/carisa/api/internal/instance/samples"

	"github.com/carisa/api/internal/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.space_test

func TestSpaceService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	s, err := space(mng)

	if assert.NoError(t, err) {
		ok, found, err := srv.Create(s)

		if assert.NoError(t, err) {
			assert.True(t, ok, "Created")
			assert.True(t, found, "Instance found")
			checkSpace(t, srv, *s)
		}
	}
}

func TestSpaceService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	inst, err := instsmpl.CreateInstance(mng)
	if err != nil {
		assert.Error(t, err, "Creating instance")
	}

	tests := []struct {
		name    string
		updated bool
		space   *Space
	}{
		{
			name:    "Creating space",
			updated: false,
			space: &Space{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				InstID: inst.ID,
			},
		},
		{
			name:    "Updating space",
			updated: true,
			space: &Space{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				InstID: inst.ID,
			},
		},
	}

	for _, tt := range tests {
		updated, found, err := srv.Put(tt.space)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Space updated"))
			assert.True(t, found, strings.Concat(tt.name, "Instance found"))
			checkSpace(t, srv, *tt.space)
		}
	}
}

func checkSpace(t *testing.T, srv Service, s Space) {
	var sr Space
	_, err := srv.Get(s.ID, &sr)
	if assert.NoError(t, err) {
		assert.Equal(t, s, sr, "Getting space")
	}
}

func TestSpaceService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	s, err := space(mng)

	if assert.NoError(t, err) {
		_, _, err := srv.Create(s)
		if assert.NoError(t, err) {
			var gets Space
			ok, err := srv.Get(s.ID, &gets)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, s, &gets, "Space returned")
			}
		}
	}
}

func TestSpaceService_ListEntes(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	link, _, err := entesmpl.CreateLinkForSpace(mng, id)

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := s.ListEntes(id, "name", tt.Ranges, 1)
			if assert.NoError(t, err) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
			}
		}
	}
}

func TestSpaceService_ListCategories(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	link, _, err := catsmpl.CreateRootLink(mng, id)

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := s.ListCategories(id, "name", tt.Ranges, 1)
			if assert.NoError(t, err) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
			}
		}
	}
}

func space(mng storage.Integration) (*Space, error) {
	inst, err := instsmpl.CreateInstance(mng)
	if err == nil {
		space := New()
		space.Name = "name"
		space.Desc = "desc"
		space.InstID = inst.ID
		return &space, nil
	}
	return nil, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	ext := srv.NewExt(cnt, crudOper.Store())
	return NewService(cnt, ext, crudOper), mng
}
