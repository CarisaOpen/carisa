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
	"github.com/carisa/api/internal/runtime"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/pkg/errors"
)

// Service implements the instance domain service
type Service struct {
	cnf   runtime.Config
	store storage.CRUD
}

// Create creates a instance into of the repository
// If the instance exists returns false
func (s *Service) Create(inst Instance) (bool, error) {
	txn := storage.NewTxn(s.store)

	txn.Find(inst.GetKey())

	create, err := s.store.Create(inst)
	if err != nil {
		return false, err
	}

	txn.DoNotFound(create)

	ctx, cancel := s.cnf.StoreWithTimeout()
	ok, err := txn.Commit(ctx)
	cancel()
	if err != nil {
		return false, errors.Wrap(err, strings.Concat("Commit creating. ", inst.ToString()))
	}

	return ok, nil
}
