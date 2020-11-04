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

	"github.com/carisa/internal/api/plugin"

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestInstance_ToString(t *testing.T) {
	i := New()
	assert.Equal(t, strings.Concat("instance: ID:", i.Key(), ", name:", i.Name), i.ToString())
}

func TestInstance_Key(t *testing.T) {
	i := New()
	assert.Equal(t, entity.ObjectKey(i.ID), i.Key())
}

func TestInstance_Nominative(t *testing.T) {
	i := Instance{}
	assert.Equal(t, entity.Descriptor{}, i.Nominative())
}

func TestInstance_ParentKey(t *testing.T) {
	i := New()
	i.SchContainer = entity.SchCategory
	assert.Equal(t, entity.Key(entity.SchCategory, i.ContainerID), i.ParentKey())
}

func TestInstance_Empty(t *testing.T) {
	i := New()
	assert.Equal(t, &Instance{}, i.Empty())
}

func TestInstance_Link(t *testing.T) {
	name := "nameinst"

	id := xid.New()
	containerID := xid.New()
	proto := Instance{
		Descriptor: entity.Descriptor{
			ID:   id,
			Name: name,
		},
		SchContainer: entity.SchCategory,
		ContainerID:  containerID,
		Category:     plugin.Query,
	}
	link := &relation.PlatformInstance{
		ID:       strings.Concat(entity.Key(entity.SchCategory, containerID), string(plugin.Query), name, entity.ObjectKey(id)),
		Name:     name,
		InstID:   id.String(),
		Category: "query",
	}

	assert.Equal(t, link, proto.Link())
}

func TestInstance_LinkName(t *testing.T) {
	i := New()
	assert.Equal(t, string(i.Category), i.LinkName())
}

func TestInstance_ReLink(t *testing.T) {
	name := "nameinstdlr"

	id := xid.New()
	containerID := xid.New()
	containerKey := entity.Key(entity.SchCategory, containerID)
	proto := Instance{
		Descriptor: entity.Descriptor{
			ID:   id,
			Name: name,
		},
		SchContainer: entity.SchCategory,
		ContainerID:  containerID,
		Category:     plugin.Query,
	}
	link := &relation.PlatformInstance{
		ID:       strings.Concat(containerKey, string(plugin.Query), name, entity.ObjectKey(id)),
		Name:     name,
		InstID:   id.String(),
		Category: "query",
	}

	dlr := storage.DLRel{
		ParentID: containerKey,
		Type:     string(plugin.Query),
	}

	assert.Equal(t, link, proto.ReLink(dlr))
}
