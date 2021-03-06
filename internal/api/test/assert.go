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

package test

import (
	"context"
	"testing"

	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/stretchr/testify/assert"
)

func CheckRelations(t *testing.T, crud storage.CRUD, name string, e storage.EntityRelation) {
	found, err := crud.Exists(context.TODO(), e.Link().Key())
	if assert.NoError(t, err, name) {
		assert.True(t, found, strings.Concat(name, "Getting link"))
	}
	found, err = crud.Exists(context.TODO(), storage.DLRKey(e.Key(), e.ParentKey()))
	if assert.NoError(t, err, name) {
		assert.True(t, found, strings.Concat(name, "Getting DLR"))
	}
}
