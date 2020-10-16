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

package convert

import (
	nethttp "net/http"
	"strconv"

	"github.com/carisa/pkg/strings"

	"github.com/carisa/pkg/http"
	"github.com/rs/xid"
)

// ParamID convert string param ID to xId type
func ParamID(c http.Context) (xid.ID, error) {
	id, err := ParamXID(c, "id")
	return id, err
}

// FilterLink gets the filter parameters.
// The parameters are entity ID, sname or gtname, top and filter type. Look at api documentation
// The default top parameter is 20
func FilterLink(c http.Context, excludeID bool) (xid.ID, string, int, bool, error) {
	id := xid.NilID()
	var err error

	if !excludeID {
		id, err = ParamXID(c, "id")
		if err != nil {
			return xid.NilID(), "", 0, false, err
		}
	}

	tops := c.QueryParam("top")
	top := 20
	if len(tops) != 0 {
		top, err = strconv.Atoi(tops)
		if err != nil {
			return xid.NilID(),
				"",
				0,
				false,
				c.HTTPError(nethttp.StatusBadRequest, "the filter top parameter has a incorrect format")
		}
	}
	if !(top >= 1 && top <= 100) {
		return xid.NilID(),
			"",
			0,
			false,
			c.HTTPError(nethttp.StatusBadRequest, "the top parameters must be between 1 and 100")
	}

	sname := c.QueryParam("sname")
	if len(sname) != 0 {
		return id, sname, top, false, nil
	}

	gtname := c.QueryParam("gtname")
	if len(gtname) != 0 {
		return id, gtname, top, true, nil
	}

	return xid.NilID(),
		"",
		0,
		false,
		c.HTTPError(nethttp.StatusBadRequest, "the filter parameters are missing. sname or qtname")
}

func ParamXID(c http.Context, name string) (xid.ID, error) {
	value := c.Param(name)

	if len(value) == 0 {
		return xid.NilID(), c.HTTPError(nethttp.StatusBadRequest, strings.Concat("the param path:", name, " not found"))
	}

	id, err := xid.FromString(value)
	if err != nil {
		return xid.NilID(), c.HTTPError(nethttp.StatusBadRequest, strings.Concat("the ", name, " has a incorrect format"))
	}
	return id, nil
}
