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

package ente

import (
	"testing"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/api/internal/entity"
	"github.com/rs/xid"

	spcsamples "github.com/carisa/api/internal/space/samples"

	"github.com/carisa/api/internal/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.ente_test

func TestEnteService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	s, err := ente(mng)

	if assert.NoError(t, err) {
		ok, found, err := srv.Create(s)

		if assert.NoError(t, err) {
			assert.True(t, ok, "Created")
			assert.True(t, found, "Space found")
			checkEnte(t, srv, *s)
		}
	}
}

func TestEnteService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	space, err := spcsamples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err, "Creating ente")
	}

	tests := []struct {
		name    string
		updated bool
		ente    *Ente
	}{
		{
			name:    "Creating ente",
			updated: false,
			ente: &Ente{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				SpaceID: space.ID,
			},
		},
		{
			name:    "Updating ente",
			updated: true,
			ente: &Ente{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				SpaceID: space.ID,
			},
		},
	}

	for _, tt := range tests {
		updated, found, err := srv.Put(tt.ente)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Ente updated"))
			assert.True(t, found, strings.Concat(tt.name, "Space found"))
			checkEnte(t, srv, *tt.ente)
		}
	}
}

func checkEnte(t *testing.T, srv Service, e Ente) {
	var er Ente
	_, err := srv.Get(e.ID, &er)
	if assert.NoError(t, err) {
		assert.Equal(t, e, er, "Getting ente")
	}
}

func TestEnteService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	s, err := ente(mng)

	if assert.NoError(t, err) {
		_, _, err := srv.Create(s)
		if assert.NoError(t, err) {
			var get Ente
			ok, err := srv.Get(s.ID, &get)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, s, &get, "Ente returned")
			}
		}
	}
}

func ente(mng storage.Integration) (*Ente, error) {
	space, err := spcsamples.CreateSpace(mng)
	if err == nil {
		ente := New()
		ente.Name = "name"
		ente.Desc = "desc"
		ente.SpaceID = space.ID
		return &ente, nil
	}
	return nil, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	return NewService(cnt, crudOper), mng
}
