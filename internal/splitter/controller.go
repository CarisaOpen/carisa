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

package splitter

import "github.com/carisa/pkg/storage"

// Controller implements the functionality when the splitter service starts, stops, etc.
type Controller struct {
	store storage.CRUD
}

// NewController builds a Controller
func NewController(data storage.CRUD) Controller {
	return Controller{
		store: data,
	}
}

func (c *Controller) Start() {

}
