/*
 *  Copyright 2019-2022 the original author or authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package factory

import (
	"testing"

	"github.com/carisa/internal/api/mock"
	"github.com/stretchr/testify/assert"
)

func TestTemplate_Build(t *testing.T) {
	sMock := mock.NewStorageFake(t)
	defer sMock.Close()

	factory := build(sMock)

	assert.NotNil(t, factory.cnt, "Container")
	assert.NotNil(t, factory.store, "Store")

	assert.NotNil(t, factory.Controller, "Controller")
}
