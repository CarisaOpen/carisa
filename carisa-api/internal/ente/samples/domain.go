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
	"github.com/carisa/api/internal/ente"
	"github.com/carisa/api/internal/mock"
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

func CreateEnte(mng storage.Integration) (ente.Ente, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	ente := ente.New()
	ente.Name = "Name"
	ente.Desc = "desc"
	_, err := crudOper.Put("loc", cnt.StoreWithTimeout, &ente)
	return ente, err
}

func CreateLinkForSpace(mng storage.Integration, spaceID xid.ID) (storage.Entity, ente.Ente, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	s := ente.New()
	s.Name = "name"
	s.Desc = "desc"
	s.SpaceID = spaceID
	link := s.Link()
	_, err := crudOper.Create("", cnt.StoreWithTimeout, link)
	return link, s, err
}

func CatLinkForEnte(mng storage.Integration, catId xid.ID, name string, enteId xid.ID) (storage.Entity, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	link := &relation.Hierarchy{
		ID:       strings.Concat(catId.String(), name, enteId.String()),
		Name:     name,
		LinkID:   enteId.String(),
		Category: false,
	}
	_, err := crudOper.Create("", cnt.StoreWithTimeout, link)
	return link, err
}

func CreateLinkProp(mng storage.Integration, enteID xid.ID) (storage.Entity, ente.Prop, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	prop := ente.NewProp()
	prop.Name = "namep"
	prop.Desc = "descp"
	prop.EnteID = enteID
	link := prop.Link()
	_, err := crudOper.Create("", cnt.StoreWithTimeout, link)
	return link, prop, err
}
