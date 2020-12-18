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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/plugin"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// Instance has n dynamic properties and them links with the associated plugin.Prototype properties
// Example: Queries. The queries has a plugin.Prototype where define the metadata of a query type
// and the Instance is the particular information of those metadata configured by users
// Is similar to classes and objects. Object will be the Instance.
type Instance struct {
	entity.Descriptor
	SchContainer string          `json:"-"`
	ContainerID  xid.ID          `json:"containerId"`
	ProtoID      xid.ID          `json:"prototypeId"`
	Category     plugin.Category `json:"-"`
}

func New() Instance {
	return Instance{
		Descriptor: entity.NewDescriptor(),
		Category:   plugin.Query,
	}
}

func (i *Instance) ToString() string {
	return strings.Concat("instance: ID:", i.Key(), ", name:", i.Name)
}

func (i *Instance) Key() string {
	return entity.ObjectKey(i.ID)
}

func (i *Instance) Nominative() entity.Descriptor {
	return i.Descriptor
}

func (i *Instance) RelName() string {
	return i.Name
}

func (i *Instance) ParentKey() string {
	return entity.Key(i.SchContainer, i.ContainerID)
}

// Link gets the link between the container and instance
func (i *Instance) Link() storage.Entity {
	return i.link(string(i.Category), i.ParentKey())
}

func (i *Instance) LinkName() string {
	return string(i.Category)
}

func (i *Instance) ReLink(dlr storage.DLRel) storage.Entity {
	return i.link(dlr.Type, dlr.ParentID)
}

func (i *Instance) link(category string, parentID string) storage.Entity {
	return &relation.PlatformInstance{
		ID:       strings.Concat(parentID, category, i.Name, i.Key()),
		Name:     i.Name,
		InstID:   i.ID.String(),
		Category: category,
	}
}

func (i *Instance) Empty() storage.EntityRelation {
	return &Instance{}
}
