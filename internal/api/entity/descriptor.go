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

package entity

import (
	"github.com/rs/xid"
)

type Descriptors interface {
	Nominative() Descriptor
}

// Descriptor describes the entity with name and description
type Descriptor struct {
	ID   xid.ID `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Desc string `json:"description,omitempty"`
}

func NewDescriptor() Descriptor {
	return Descriptor{
		ID: xid.New(),
	}
}

func (d *Descriptor) AutoID() {
	d.ID = xid.New()
}
