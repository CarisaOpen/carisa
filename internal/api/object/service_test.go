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

package object

import (
	"testing"

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"

	psamples "github.com/carisa/internal/api/plugin/samples"

	"github.com/carisa/internal/api/plugin"

	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/samples"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/internal/api/test"
	"github.com/carisa/pkg/storage"

	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.object_test

func TestObjectService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	protoID := xid.New()
	proto, err := psamples.CreatePlugin(mng, plugin.Query, protoID)
	if err != nil {
		assert.NoError(t, err, "Creating plugin")
		return
	}

	inst, err := instance(mng, proto)
	if err != nil {
		assert.NoError(t, err, "Creating instance")
		return
	}

	ok, foundp, foundc, err := srv.Create(inst)

	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
		assert.True(t, foundp, "Checking plugin prototype")
		assert.True(t, foundc, "Checking container")
		checkInst(t, srv, "Checking relations", *inst)
	}
}

func TestObjectService_Put(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	protoID := xid.New()
	proto, err := psamples.CreatePlugin(mng, plugin.Query, protoID)
	if err != nil {
		assert.NoError(t, err, "Creating plugin")
		return
	}
	container, err := samples.CreateEntityMock(mng, entity.SchCategory)
	if err != nil {
		assert.NoError(t, err, "Creating container")
		return
	}

	tests := []struct {
		name    string
		updated bool
		inst    *Instance
	}{
		{
			name:    "Creating Instance.",
			updated: false,
			inst: &Instance{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				ProtoID:      proto.ID,
				SchContainer: entity.SchCategory,
				ContainerID:  container.ID,
			},
		},
		{
			name:    "Updating Instance.",
			updated: true,
			inst: &Instance{
				Descriptor: entity.Descriptor{
					ID:   xid.NilID(),
					Name: "name",
					Desc: "desc",
				},
				SchContainer: entity.SchCategory,
				ContainerID:  container.ID,
			},
		},
	}

	for _, tt := range tests {
		updated, foundp, foundc, err := srv.Put(tt.inst)
		if assert.NoError(t, err) {
			assert.Equal(t, updated, tt.updated, strings.Concat(tt.name, "Instance updated"))
			assert.True(t, foundp, strings.Concat(tt.name, "Seeking plugin"))
			assert.True(t, foundc, strings.Concat(tt.name, "Seeking container"))
			checkInst(t, srv, tt.name, *tt.inst)
		}
	}
}

func checkInst(t *testing.T, srv Service, name string, inst Instance) {
	var ir Instance
	_, err := srv.Get(inst.ID, &ir)
	if assert.NoError(t, err) {
		assert.Equal(t, inst, ir, "Getting instance")
	}
	test.CheckRelations(t, srv.crud.Store(), name, &inst)
}

func TestObjectService_Get(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	protoID := xid.New()
	proto, err := psamples.CreatePlugin(mng, plugin.Query, protoID)
	if err != nil {
		assert.NoError(t, err, "Creating plugin")
		return
	}

	inst, err := instance(mng, proto)

	if assert.NoError(t, err) {
		_, _, _, err := srv.Create(inst)
		if assert.NoError(t, err) {
			var get Instance
			ok, err := srv.Get(inst.ID, &get)
			if assert.NoError(t, err) {
				assert.True(t, ok, "Get ok")
				assert.Equal(t, inst, &get, "Instance returned")
			}
		}
	}
}

func TestObjectService_ListQueries(t *testing.T) {
	tests := samples.TestList()

	s, mng := newServiceFaked(t)
	defer mng.Close()

	id := xid.New()
	inst := New()
	inst.SchContainer = entity.SchCategory
	inst.ContainerID = id
	inst.Category = plugin.Query
	inst.Name = "namei"
	link := inst.Link()

	_, err := s.crud.Create("", s.cnt.StoreWithTimeout, link)
	if err != nil {
		assert.Error(t, err, "Create query link")
		return
	}

	for _, tt := range tests {
		list, err := s.ListInstances(inst.SchContainer, id, plugin.Query, "namei", tt.Ranges, 2)
		if assert.NoError(t, err, tt.Name) {
			assert.Equalf(t, link, list[0], "Ranges: %v", tt.Name)
		}
	}
}

func instance(mng storage.Integration, proto plugin.Prototype) (*Instance, error) {
	entityMock, err := samples.CreateEntityMock(mng, entity.SchCategory)
	if err != nil {
		return nil, err
	}
	inst := New()
	inst.Name = "name"
	inst.Desc = "desc"
	inst.SchContainer = entity.SchCategory
	inst.ContainerID = entityMock.ID
	inst.ProtoID = proto.ID
	return &inst, nil
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	ext := service.NewExt(cnt, crudOper.Store())
	protos := plugin.NewService(cnt, ext, crudOper)
	return NewService(cnt, ext, crudOper, &protos), mng
}
