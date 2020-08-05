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

	"github.com/carisa/api/internal/entity"
	"github.com/rs/xid"

	"github.com/carisa/api/internal/samples"

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
		}
	}
}

func TestInstanceService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	inst, err := samples.CreateInstance(mng)
	if err != nil {
		assert.Error(t, err, "Creating instance")
	}

	tests := []struct {
		updated bool
		space   *Space
	}{
		{
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
			assert.Equal(t, updated, tt.updated, "Space updated")
			assert.True(t, found, "Instance found")
		}
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
	inst, err := samples.CreateInstance(mng)
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
