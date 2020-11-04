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
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
)

type EntityMock struct {
	entity.Descriptor
	scheme string
}

func NewEM() EntityMock {
	return EntityMock{
		Descriptor: entity.NewDescriptor(),
	}
}

func (e EntityMock) ToString() string {
	return strings.Concat("EntityMock: ", e.Name)
}

func (e EntityMock) Key() string {
	return strings.Concat(e.scheme, e.ID.String())
}

func CreateEntityMock(mng storage.Integration, scheme string) (EntityMock, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	e := NewEM()
	e.scheme = scheme
	e.Name = "name"
	e.Desc = "desc"
	_, err := crudOper.Put("loc", cnt.StoreWithTimeout, &e)
	return e, err
}
