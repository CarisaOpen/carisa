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

package validator

import (
	nethttp "net/http"

	"github.com/carisa/internal/api/entity"
	"github.com/carisa/pkg/http"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

// Descriptor validates that the name or description is not empty
func Descriptor(ctx http.Context, d entity.Descriptor) error {
	if err := ctx.NoEmpty("name", d.Name); err != nil {
		return err
	}
	if err := ctx.NoEmpty("description", d.Desc); err != nil {
		return err
	}
	if err := ctx.MaxLen("name", d.Name, 50); err != nil {
		return err
	}
	if err := ctx.MaxLen("description", d.Desc, 500); err != nil {
		return err
	}
	return nil
}

func ID(ctx http.Context, id xid.ID) error {
	if id.IsNil() {
		return ctx.HTTPError(nethttp.StatusBadRequest, strings.Concat("the property: 'ID' can not be empty"))
	}
	return nil
}
