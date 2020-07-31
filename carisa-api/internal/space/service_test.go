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
	"testing"

	"github.com/carisa/api/internal/samples"

	"github.com/carisa/api/internal/mock"
	"github.com/carisa/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Verify the crud integration. For all rest test look at http.handler.space_test

func TestInstanceService_Create(t *testing.T) {
	srv, mng := newServiceFaked(t)
	defer mng.Close()

	s, err := space(mng)

	if assert.NoError(t, err) {
		ok, found, err := srv.Create(s)

		if assert.NoError(t, err) {
			assert.True(t, ok, "Created")
			assert.True(t, found, "Instance found")
		}
	}
}

func space(mng storage.Integration) (*Space, error) {
	inst, err := samples.CreateInstance(mng)
	if err == nil {
		space := NewSpace()
		space.Name = "name"
		space.Desc = "desc"
		space.InstID = inst.ID
		return &space, nil
	}
	return nil, err
}

func newServiceFaked(t *testing.T) (Service, storage.Integration) {
	mng, cnt, crudOper := mock.NewFullCrudOperFaked(t)
	return NewService(cnt, crudOper), mng
}
