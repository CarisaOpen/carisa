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

	"github.com/rs/xid"

	"github.com/carisa/pkg/storage"
)

func TestList() []struct {
	Name   string
	Ranges bool
} {
	return []struct {
		Name   string
		Ranges bool
	}{
		{
			Name:   "Find by start name",
			Ranges: true,
		},
		{
			Name:   "Find by range",
			Ranges: false,
		},
	}
}

func TestCreateWithError(method string) []struct {
	Name     string
	Body     string
	MockOper func(txn *storage.ErrMockCRUDOper)
	Status   int
} {
	return []struct {
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
			Body:   `{"Name":"","description":"desc"}`,
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:     "Creating the entity. Error creating",
			Body:     `{"Name":"Name","description":"desc"}`,
			MockOper: func(s *storage.ErrMockCRUDOper) { s.Activate(method) },
			Status:   nethttp.StatusInternalServerError,
		},
	}
}

func TestPutWithError(method string, params map[string]string) []struct {
	Name     string
	Params   map[string]string
	Body     string
	MockOper func(txn *storage.ErrMockCRUDOper)
	Status   int
} {
	return []struct {
		Name     string
		Params   map[string]string
		Body     string
		MockOper func(txn *storage.ErrMockCRUDOper)
		Status   int
	}{
		{
			Name:   "Body wrong. Bad request",
			Params: params,
			Body:   "{df",
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:   "Getting the identifier (ID) param. Bad request",
			Params: map[string]string{"i": ""},
			Body:   `{"Name":"Name","description":"desc"}`,
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:   "Descriptor validation. Bad request",
			Params: params,
			Body:   `{"Name":"Name","description":""}`,
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:     "Putting the entity. Error putting",
			Params:   params,
			Body:     `{"Name":"Name","description":"desc"}`,
			MockOper: func(s *storage.ErrMockCRUDOper) { s.Activate(method) },
			Status:   nethttp.StatusInternalServerError,
		},
	}
}

func TestGetWithError() []struct {
	Name     string
	Param    map[string]string
	MockOper func(txn *storage.ErrMockCRUDOper)
	Status   int
} {
	return []struct {
		Name     string
		Param    map[string]string
		MockOper func(txn *storage.ErrMockCRUDOper)
		Status   int
	}{
		{
			Name:   "Param not found. Bad request",
			Param:  map[string]string{"i": ""},
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:     "Get error. Internal server error",
			Param:    map[string]string{"id": xid.NilID().String()},
			MockOper: func(s *storage.ErrMockCRUDOper) { s.Store().(*storage.ErrMockCRUD).Activate("Get") },
			Status:   nethttp.StatusInternalServerError,
		},
	}
}

func TestListError() []struct {
	Name     string
	Param    map[string]string
	QParam   map[string]string
	MockOper func(txn *storage.ErrMockCRUDOper)
	Status   int
} {
	return []struct {
		Name     string
		Param    map[string]string
		QParam   map[string]string
		MockOper func(txn *storage.ErrMockCRUDOper)
		Status   int
	}{
		{
			Name:   "Param not found. Bad request",
			Param:  map[string]string{"i": ""},
			Status: nethttp.StatusBadRequest,
		},
		{
			Name:     "List error. Internal server error",
			Param:    map[string]string{"id": xid.NilID().String()},
			QParam:   map[string]string{"sname": "sname"},
			MockOper: func(s *storage.ErrMockCRUDOper) { s.Store().(*storage.ErrMockCRUD).Activate("StartKey") },
			Status:   nethttp.StatusInternalServerError,
		},
	}
}
