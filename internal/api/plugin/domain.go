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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
)

// Category of the plugin
type Category string

const (
	Query Category = "query"
)

func Plugins() []Category {
	return []Category{Query}
}

// It allows make prototype of dynamic object.
// Each object.Instance have n dynamic properties.
// Each object.Instance define the representative name.
// The plugin developers could define the necessary information what the plugin needs
// with the Prototype and PrototypeProperty.
// The user interface will be built dynamically with this information.
// External users could make dynamic object object.Instance using plugin Prototype.
// See object.Instance
type Prototype struct {
	entity.Descriptor
	Category Category `json:"-"`
}

func New() Prototype {
	return Prototype{
		Descriptor: entity.NewDescriptor(),
		Category:   Query,
	}
}

func (p *Prototype) ToString() string {
	return strings.Concat("plugin: ID:", p.Key(), ", name:", p.Name)
}

func (p *Prototype) Key() string {
	return entity.PluginKey(p.ID)
}

func (p *Prototype) Nominative() entity.Descriptor {
	return p.Descriptor
}

func (p *Prototype) RelName() string {
	return p.Name
}

// ParentKey gets the Virtual ID. This parent is not stored in the DB
func (p *Prototype) ParentKey() string {
	return storage.Virtual
}

// Link gets the link between the platform and plugin
func (p *Prototype) Link() storage.Entity {
	return p.link(string(p.Category), p.ParentKey())
}

func (p *Prototype) LinkName() string {
	return string(p.Category)
}

func (p *Prototype) ReLink(dlr storage.DLRel) storage.Entity {
	return p.link(dlr.Type, dlr.ParentID)
}

func (p *Prototype) link(category string, parentID string) storage.Entity {
	return &relation.PlatformPlugin{
		ID:       strings.Concat(parentID, category, p.Name, p.Key()),
		Name:     p.Name,
		ProtoID:  p.ID.String(),
		Category: category,
	}
}

func (p *Prototype) Empty() storage.EntityRelation {
	return &Prototype{}
}
