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
	"github.com/rs/xid"
)

// Instance has n dynamic properties
// External users could instance a dynamic object using dynamic object instance
// Example: Queries. The queries has a object prototype where define the metadata of a query type
// and the instance is the particular information of those metadata configured by users
// Is similar to classes and objects. Object will be the instance.
type Instance struct {
	entity.Descriptor
	PrototypeID xid.ID `json:"prototypeID"`
}

func NewInstance() Instance {
	return Instance{
		Descriptor: entity.NewDescriptor(),
	}
}
