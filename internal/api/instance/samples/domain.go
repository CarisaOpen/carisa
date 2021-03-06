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
	"github.com/carisa/internal/api/instance"
	"github.com/carisa/internal/api/mock"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/storage"
)

func CreateInstance(mng storage.Integration) (instance.Instance, error) {
	cnt, crudOper := mock.NewCrudOperFaked(mng)
	ext := service.NewExt(cnt, crudOper.Store())
	srv := instance.NewService(cnt, ext, crudOper)
	inst := instance.New()
	inst.Name = "Name"
	inst.Desc = "desc"
	_, err := srv.Create(&inst)
	return inst, err
}
