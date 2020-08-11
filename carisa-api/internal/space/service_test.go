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

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/entity"
	"github.com/rs/xid"

	instsamples "github.com/carisa/api/internal/instance/samples"

	"github.com/carisa/api/internal/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.space_test

func TestInstanceService_Create(t *testing.T) {
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

func TestInstanceService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	inst, err := instsamples.CreateInstance(mng)
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

func TestInstanceService_Get(t *testing.T) {
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

func space(mng storage.Integration) (*Space, error) {
	inst, err := instsamples.CreateInstance(mng)
	if err == nil {
		space := NewSpace()
		space.Name = "name"
		space.Desc = "desc"
		space.InstID = inst.ID
		return &space, nil
	}
	return nil, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	return NewService(cnt, crudOper), mng
}
