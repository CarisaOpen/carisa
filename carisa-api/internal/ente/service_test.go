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
	"context"
	"testing"

	"github.com/carisa/api/internal/test"

	"github.com/carisa/api/internal/service"

	"github.com/carisa/api/internal/samples"

	"github.com/carisa/api/internal/runtime"

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
			checkEnte(t, srv, "Checking relations", *s)
		}
	}
}

func TestEnteService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	space, err := spcsamples.CreateSpace(mng)
	if err != nil {
		assert.Error(t, err, "Creating space")
	}

	tests := []struct {
		name    string
		updated bool
		ente    *Ente
	}{
		{
			name:    "Creating prop",
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
			name:    "Updating prop",
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
			checkEnte(t, srv, tt.name, *tt.ente)
		}
	}
}

func checkEnte(t *testing.T, srv Service, name string, e Ente) {
	var er Ente
	_, err := srv.Get(e.ID, &er)
	if assert.NoError(t, err) {
		assert.Equal(t, e, er, "Getting prop")
	}
	test.CheckRelations(t, srv.crud.Store(), name, &e)
}

func TestEnteService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	e, err := ente(mng)

	if assert.NoError(t, err) {
		_, _, err := srv.Create(e)
		if assert.NoError(t, err) {
			var get Ente
			ok, err := srv.Get(e.ID, &get)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, e, &get, "Ente returned")
			}
		}
	}
}

func TestEnteService_LinkToCat(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	ente, err := createEnte(srv.cnt, srv.crud)
	if err != nil {
		assert.Error(t, err, "Creating prop")
	}

	cat, err := samples.CreateEntityMock(mng)
	if err != nil {
		assert.Error(t, err, "Creating category")
	}

	sfound, tfound, rel, err := srv.LinkToCat(ente.ID, cat.ID)

	if assert.NoError(t, err) {
		assert.True(t, sfound, "Ente found")
		assert.True(t, tfound, "Category found")

		found, err := srv.crud.Store().Exists(context.TODO(), rel.Key())
		if assert.NoError(t, err) {
			assert.True(t, found, "Getting link")
		}
		found, err = srv.crud.Store().Exists(context.TODO(), storage.DLRKey(ente.Key(), cat.Key()))
		if assert.NoError(t, err) {
			assert.True(t, found, "Getting DLR")
		}
	}
}

func TestEnteService_ListProps(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	prop := NewProp()
	prop.EnteID = id
	prop.Name = "namep"
	link := prop.Link()

	_, err := s.crud.Create("", s.cnt.StoreWithTimeout, link)

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := s.ListProps(id, "namep", tt.Ranges, 1)
			if assert.NoError(t, err, tt.Name) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
			}
		}
	}
}

func TestEnteService_CreateProp(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	prop, err := prop(srv.cnt, srv.crud)

	if assert.NoError(t, err) {
		ok, found, err := srv.CreateProp(prop)

		if assert.NoError(t, err) {
			assert.True(t, ok, "Created")
			assert.True(t, found, "Ente found")
			checkProp(t, srv, "Checking relations", *prop)
		}
	}
}

func TestEnteService_PutProp(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	ente, err := createEnte(srv.cnt, srv.crud)
	if err != nil {
		assert.Error(t, err, "Creating prop")
	}

	tests := []struct {
		name    string
		updated bool
		prop    *Prop
	}{
		{
			name:    "Creating property",
			updated: false,
			prop: &Prop{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				EnteID: ente.ID,
			},
		},
		{
			name:    "Updating property",
			updated: true,
			prop: &Prop{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				Type:   entity.Boolean,
				EnteID: ente.ID,
			},
		},
	}

	for _, tt := range tests {
		updated, found, err := srv.PutProp(tt.prop)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Property updated"))
			assert.True(t, found, strings.Concat(tt.name, "Ente found"))
			checkProp(t, srv, tt.name, *tt.prop)
		}
	}
}

func checkProp(t *testing.T, srv Service, name string, p Prop) {
	var prop Prop
	_, err := srv.GetProp(p.ID, &prop)
	if assert.NoError(t, err) {
		assert.Equal(t, p, prop, "Getting property")
	}
	test.CheckRelations(t, srv.crud.Store(), name, &p)
}

func TestEnteService_GetProp(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	prop, err := prop(srv.cnt, srv.crud)

	if assert.NoError(t, err) {
		_, _, err := srv.CreateProp(prop)
		if assert.NoError(t, err) {
			var get Prop
			ok, err := srv.GetProp(prop.ID, &get)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, prop, &get, "Property returned")
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

func prop(cnt *runtime.Container, crud storage.CrudOperation) (*Prop, error) {
	ente, err := createEnte(cnt, crud)
	if err == nil {
		prop := NewProp()
		prop.Name = "name"
		prop.Desc = "desc"
		prop.Type = entity.Decimal
		prop.EnteID = ente.ID
		return &prop, nil
	}
	return nil, err
}

func createEnte(cnt *runtime.Container, crud storage.CrudOperation) (Ente, error) {
	ente := New()
	_, err := crud.Put("loc", cnt.StoreWithTimeout, &ente)
	return ente, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crudOper.Store())
	return NewService(cnt, ext, crudOper), mng
}
