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

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestPrototype_ToString(t *testing.T) {
	p := New()
	assert.Equal(t, strings.Concat("plugin: ID:", p.Key(), ", name:", p.Name), p.ToString())
}

func TestPrototype_Key(t *testing.T) {
	p := New()
	assert.Equal(t, p.ID.String(), p.Key())
}

func TestPrototype_Nominative(t *testing.T) {
	p := Prototype{}
	assert.Equal(t, entity.Descriptor{}, p.Nominative())
}

func TestPrototype_ParentKey(t *testing.T) {
	p := New()
	assert.Equal(t, storage.Virtual, p.ParentKey())
}

func TestPrototype_Empty(t *testing.T) {
	p := New()
	assert.Equal(t, &Prototype{}, p.Empty())
}

func TestPrototype_Link(t *testing.T) {
	name := "nameproto"

	id := xid.New()
	proto := Prototype{
		Descriptor: entity.Descriptor{
			ID:   id,
			Name: name,
		},
		Category: Query,
	}
	link := &relation.PlatformPlugin{
		ID:       strings.Concat(storage.Virtual, string(Query), name, id.String()),
		Name:     name,
		ProtoID:  id.String(),
		Category: "query",
	}

	assert.Equal(t, link, proto.Link())
}

func TestPrototype_LinkName(t *testing.T) {
	proto := New()
	assert.Equal(t, string(proto.Category), proto.LinkName())
}

func TestPrototype_ReLink(t *testing.T) {
	name := "nameprotodlr"

	id := xid.New()
	proto := Prototype{
		Descriptor: entity.Descriptor{
			ID:   id,
			Name: name,
		},
		Category: Query,
	}
	link := &relation.PlatformPlugin{
		ID:       strings.Concat(storage.Virtual, string(Query), name, id.String()),
		Name:     name,
		ProtoID:  id.String(),
		Category: "query",
	}

	dlr := storage.DLRel{
		ParentID: storage.Virtual,
		Type:     string(Query),
	}

	assert.Equal(t, link, proto.ReLink(dlr))
}
