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

// As in a key value database there is not scheme, this package generates keys keep mind the scheme
package entity

import (
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

const (
	SchInstance string = "I"
	SchSpace    string = "S"
	SchEnte     string = "E"
	SchEnteProp string = "EP"
	SchCategory string = "C"
	SchCatProp  string = "CP"
	SchPlugin   string = "P"
	SchObject   string = "O"
)

func Key(scheme string, id xid.ID) string {
	return strings.Concat(scheme, id.String())
}

func InstKey(id xid.ID) string {
	return Key(SchInstance, id)
}

func SpaceKey(id xid.ID) string {
	return Key(SchSpace, id)
}

func EnteKey(id xid.ID) string {
	return Key(SchEnte, id)
}

func EntePropKey(id xid.ID) string {
	return Key(SchEnteProp, id)
}

func CategoryKey(id xid.ID) string {
	return Key(SchCategory, id)
}

func CatPropKey(id xid.ID) string {
	return Key(SchCatProp, id)
}

func PluginKey(id xid.ID) string {
	return Key(SchPlugin, id)
}

func ObjectKey(id xid.ID) string {
	return Key(SchObject, id)
}
