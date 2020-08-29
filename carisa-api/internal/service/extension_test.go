/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and  limitations under the License.
 */

package service

import (
	"github.com/carisa/api/internal/mock"
	"github.com/carisa/api/internal/relation"
	"github.com/carisa/api/internal/relation/samples"
	"github.com/carisa/pkg/storage"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtension_List(t *testing.T) {
	tests := []struct {
		name   string
		ranges bool
	}{
		{
			name:   "Find by start name",
			ranges: true,
		},
		{
			name:   "Find by range",
			ranges: false,
		},
	}

	mng, cnt, oper := mock.NewFullCrudOperFaked(t)
	defer mng.Close()

	id := xid.New()
	link, _, err := samples.CreateSpaceLink(mng, id)

	ext := Extension{
		cnt:  cnt,
		crud: oper.Store(),
	}
	empty := func() storage.Entity { return &relation.InstSpace{} }

	if assert.NoError(t, err) {
		for _, tt := range tests {
			list, err := ext.List(id, "name", tt.ranges, 1, empty)
			if assert.NoError(t, err) {
				assert.Equalf(t, link, list[0], "Ranges: %v", tt.name)
			}
		}
	}
}
