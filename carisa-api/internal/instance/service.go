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

package instance

import (
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
)

const locService = "instance.service"

// Service implements CRUD operations for the instance domain service
type Service struct {
	cnt   runtime.Container
	store storage.CRUD
}

// NewService builds a instance service
func NewService(cnt runtime.Container, store storage.CRUD) Service {
	return Service{
		cnt:   cnt,
		store: store,
	}
}

// Create creates a instance into of the repository
// If the instance exists returns false
func (s *Service) Create(inst *Instance) (bool, error) {
	inst.AutoID()

	txn := s.cnt.NewTxn(s.store)
	txn.Find(inst.GetKey())

	create, err := s.store.Create(inst)
	if err != nil {
		s.cnt.Log.ErrorE(err, locService)
		return false, err
	}

	txn.DoNotFound(create)

	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, s.cnt.Log.ErrWrap(err, "commit creating", locService, logging.String("instance", inst.ToString()))
	}

	return ok, nil
}
