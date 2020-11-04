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

package category

import (
	"testing"

	"github.com/carisa/internal/api/test"

	"github.com/carisa/internal/api/ente"

	"github.com/carisa/internal/api/samples"

	"github.com/carisa/internal/api/entity"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"

	"github.com/rs/xid"

	spcsamples "github.com/carisa/internal/api/space/samples"

	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.category_test

func TestCatService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	tests := []struct {
		name string
		root bool
	}{
		{
			name: "Creating category into space.",
			root: true,
		},
		{
			name: "Creating category into other category.",
			root: false,
		},
	}

	for _, tt := range tests {
		cat, err := category(mng, &srv, tt.root)
		if assert.NoError(t, err) {
			ok, found, err := srv.Create(cat)

			if assert.NoError(t, err, tt.name) {
				assert.True(t, ok, strings.Concat(tt.name, "Created"))
				assert.True(t, found, strings.Concat(tt.name, "Parent found"))
				checkCat(t, tt.name, srv, *cat)
			}
		}
	}
}

func TestCatService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	space, err := spcsamples.CreateSpace(mng)
	if err != nil {
		assert.NoError(t, err, "Creating space")
		return
	}

	tests := []struct {
		name    string
		updated bool
		cat     *Category
	}{
		{
			name:    "Creating category",
			updated: false,
			cat: &Category{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				ParentID: space.ID,
				Root:     true,
			},
		},
		{
			name:    "Updating category",
			updated: true,
			cat: &Category{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name1",
					Desc: "desc1",
				},
				ParentID: space.ID,
				Root:     true,
			},
		},
	}

	for _, tt := range tests {
		updated, found, err := srv.Put(tt.cat)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "CAtegory updated"))
			assert.True(t, found, strings.Concat(tt.name, "Space found"))
			checkCat(t, tt.name, srv, *tt.cat)
		}
	}
}

func checkCat(t *testing.T, name string, srv Service, cat Category) {
	var catr Category
	_, err := srv.Get(cat.ID, &catr)
	if assert.NoError(t, err, name) {
		assert.Equal(t, cat, catr, strings.Concat(name, "Getting category"))
	}
	test.CheckRelations(t, srv.crud.Store(), name, &cat)
}

func TestCatService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	cat, err := category(mng, &srv, true)

	if assert.NoError(t, err) {
		_, _, err := srv.Create(cat)
		if assert.NoError(t, err) {
			var get Category
			ok, err := srv.Get(cat.ID, &get)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, cat, &get, "Category returned")
			}
		}
	}
}

func TestCatService_ListCategories(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	cat := New()
	cat.ParentID = id
	cat.Name = "namep"
	link := cat.Link()

	_, err := s.crud.Create("", s.cnt.StoreWithTimeout, link)
	if err != nil {
		assert.NoError(t, err, "Create category link to category")
		return
	}

	for _, tt := range tests {
		list, err := s.ListCategories(id, "namep", tt.Ranges, 2)
		if assert.NoError(t, err, tt.Name) {
			assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
		}
	}
}

func TestCatService_ListProps(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	cat := NewProp()
	cat.CatID = id
	cat.Name = "namep"
	link := cat.Link()

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

func TestCatService_CreateProp(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	prop, err := prop(srv.cnt, srv.crud)

	if assert.NoError(t, err) {
		ok, found, err := srv.CreateProp(prop)

		if assert.NoError(t, err) {
			assert.True(t, ok, "Created")
			assert.True(t, found, "Category found")
			checkProp(t, srv, "Checking relations", *prop)
		}
	}
}

func TestCatService_PutProp(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()
	cat, err := createCat(srv.cnt, srv.crud)
	if err != nil {
		assert.NoError(t, err, "Creating category")
		return
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
				CatID: cat.ID,
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
				Type:  entity.Boolean,
				CatID: cat.ID,
			},
		},
	}

	for _, tt := range tests {
		updated, found, err := srv.PutProp(tt.prop)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Property updated"))
			assert.True(t, found, strings.Concat(tt.name, "Category found"))
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

func TestCatService_GetProp(t *testing.T) {
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

func category(mng storage.Integration, srv *Service, root bool) (*Category, error) {
	var id xid.ID
	if root {
		space, err := spcsamples.CreateSpace(mng)
		if err != nil {
			return nil, err
		}
		id = space.ID
	} else {
		cat, err := createCat(srv.cnt, srv.crud)
		if err != nil {
			return nil, err
		}
		id = cat.ID
	}
	categ := New()
	categ.Name = "name"
	categ.Desc = "desc"
	categ.ParentID = id
	categ.Root = root
	return &categ, nil
}

func prop(cnt *runtime.Container, crud storage.CrudOperation) (*Prop, error) {
	cat, err := createCat(cnt, crud)
	if err == nil {
		prop := NewProp()
		prop.Name = "name"
		prop.Desc = "desc"
		prop.Type = entity.Decimal
		prop.CatID = cat.ID
		return &prop, nil
	}
	return nil, err
}

func createCat(cnt *runtime.Container, crud storage.CrudOperation) (Category, error) {
	cat := New()
	_, err := crud.Put("loc", cnt.StoreWithTimeout, &cat)
	return cat, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crudOper.Store())
	entesrv := ente.NewService(cnt, ext, crudOper)
	return NewService(cnt, ext, crudOper, &entesrv), mng
}
