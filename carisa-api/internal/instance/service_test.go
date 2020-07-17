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

package instance

import (
	"testing"

	"github.com/carisa/api/internal/entity"
	"github.com/carisa/api/internal/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all test look at http.handler.instance_test
func TestInstanceService_Create(t *testing.T) {
	i := Instance{
		Descriptor: entity.Descriptor{
			Name: "name",
			Desc: "desc",
		},
	}
	s, mng := newServiceFaked(t)
	defer mng.Close()

	ok, err := s.Create(&i)

	if assert.NoError(t, err) {
		assert.True(t, ok, "Created")
	}
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng := mock.NewStorageFake(t)
	cnt := mock.NewContainerFake()
	crud := storage.NewCrudOperation(mng.Store(), cnt.Log, storage.NewTxn)
	return NewService(cnt, crud), mng
}
