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
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/plugin"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
)

func CreatePlugin(mng storage.Integration, cat plugin.Category, id xid.ID) (plugin.Prototype, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	proto := plugin.New()
	proto.ID = id
	proto.Category = cat
	proto.Name = "nameproto"
	proto.Desc = "descproto"
	_, err := crudOper.Create("", cnt.StoreWithTimeout, &proto)
	return proto, err
}

func CreateLinkPlugin(mng storage.Integration, cat plugin.Category) (storage.Entity, plugin.Prototype, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	proto := plugin.New()
	proto.Category = cat
	proto.Name = "nameproto"
	proto.Desc = "descproto"
	link := proto.Link()
	_, err := crudOper.Create("", cnt.StoreWithTimeout, link)
	return link, proto, err
}
