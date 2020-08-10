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

package samples

import (
	"github.com/carisa/api/internal/mock"
	"github.com/carisa/api/internal/space"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
)

func CreateSpaceLink(mng storage.Integration, instanceID xid.ID) (storage.Entity, space.Space, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	s := space.NewSpace()
	s.Name = "name"
	s.Desc = "desc"
	s.InstID = instanceID
	link := s.Link()
	_, err := crudOper.Create("", cnt.StoreWithTimeout, link)
	return link, s, err
}
