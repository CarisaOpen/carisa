/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package instance

import (
	"github.com/carisa/api/internal/common"
	"github.com/carisa/pkg/strings"
)

// Instance represents a set of spaces. Each space can have several dashboard.
// Each instance is independently of another instance in all system
type Instance struct {
	id string
	common.Descriptor
}

func (i Instance) ToString() string {
	return strings.Concat("Instance:", i.id)
}

func (i Instance) GetKey() string {
	return i.id
}