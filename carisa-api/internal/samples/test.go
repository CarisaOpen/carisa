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

package samples

import (
	nethttp "net/http"

	"github.com/carisa/pkg/storage"
)

func TestCreateWithError(method string) []struct {
	Name     string
	Body     string
	MockOper func(txn *storage.ErrMockCRUDOper)
	Status   int
} {
	tests := []struct {
		Name     string
		Body     string
		MockOper func(txn *storage.ErrMockCRUDOper)
		Status   int
	}{
		{
			Name:   "Body wrong. Bad request",
			Body:   "{df",
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:   "Descriptor validation. Bad request",
			Body:   `{"name":"","description":"desc"}`,
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:     "Creating the entity. Error creating",
			Body:     `{"name":"name","description":"desc"}`,
			MockOper: func(s *storage.ErrMockCRUDOper) { s.Activate(method) },
			Status:   nethttp.StatusInternalServerError,
		},
	}
	return tests
}
