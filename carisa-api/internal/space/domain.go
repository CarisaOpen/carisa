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
	"github.com/carisa/api/internal/entity"
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// The space splits the instance in logic categories.
// Each space can have several entes, dashboard, etc...
type Space struct {
	entity.Descriptor
	InstID xid.ID `json:"instanceId"` // Instance container
}

func New() Space {
	return Space{
		Descriptor: entity.NewDescriptor(),
	}
}

func (s *Space) ToString() string {
	return strings.Concat("space: ID:", s.Key(), ", name:", s.Name)
}

func (s *Space) Key() string {
	return s.ID.String()
}

func (s *Space) Nominative() entity.Descriptor {
	return s.Descriptor
}

func (s *Space) RelKey() string {
	return strings.Concat(s.InstID.String(), s.Name, s.Key())
}

func (s *Space) RelName() string {
	return s.Name
}

// ParentKey gets the instance ID
func (s *Space) ParentKey() string {
	return s.InstID.String()
}

func (s *Space) SetParentKey(value string) error {
	id, err := xid.FromString(value)
	if err != nil {
		return err
	}
	s.InstID = id
	return nil
}

// Link gets the link between instance and space
func (s *Space) Link() storage.Entity {
	return &relation.InstSpace{
		ID:      s.RelKey(),
		Name:    s.Name,
		SpaceID: s.ID.String(),
	}
}

func (s *Space) Empty() storage.EntityRelation {
	return &Space{}
}
