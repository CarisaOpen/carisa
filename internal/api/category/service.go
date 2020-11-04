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

package category

import (
	"github.com/carisa/internal/api/ente"
	"github.com/carisa/internal/api/entity"
	"github.com/carisa/internal/api/relation"
	"github.com/carisa/internal/api/runtime"
	"github.com/carisa/internal/api/service"
	"github.com/carisa/pkg/logging"
	"github.com/carisa/pkg/storage"
	"github.com/carisa/pkg/strings"
	"github.com/rs/xid"
)

const locService = "category.service"

// Service implements CRUD operations for the category
type Service struct {
	cnt     *runtime.Container
	ext     *service.Extension
	crud    storage.CrudOperation
	entesrv *ente.Service
}

// NewService builds a category service
func NewService(
	cnt *runtime.Container,
	ext *service.Extension,
	crud storage.CrudOperation,
	entesrv *ente.Service) Service {
	//
	return Service{
		cnt:     cnt,
		ext:     ext,
		crud:    crud,
		entesrv: entesrv,
	}
}

// Create creates a category into of the repository and links category and space or other category.
// If the category exists return false in the first param returned.
// If the space or category doesn't exist return false in the second param returned.
func (s *Service) Create(cat *Category) (bool, bool, error) {
	cat.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, cat)
}

// Put creates or updates a category into of the repository.
// If the category exists return true in the first param returned otherwise return false.
// If the space or cat doesn't exist return false in the second param returned.
func (s *Service) Put(cat *Category) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, cat)
}

// Get gets the category from storage
func (s *Service) Get(id xid.ID, cat *Category) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, entity.CategoryKey(id), cat)
	cancel()
	return ok, err
}

// ListCategories lists categories depending of 'ranges' parameter.
// Look at service.List
func (s *Service) ListCategories(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(entity.CategoryKey(id), name, ranges, top, func() storage.Entity { return &relation.Hierarchy{} })
}

// ListProps lists properties depending ranges parameter.
// Look at service.List
func (s *Service) ListProps(id xid.ID, name string, ranges bool, top int) ([]storage.Entity, error) {
	return s.ext.List(
		entity.CategoryKey(id),
		strings.Concat(relation.CatPropLn, name),
		ranges,
		top,
		func() storage.Entity { return &relation.CategoryProp{} })
}

// CreateProp creates a property into of the repository and links category property and category.
// If the property exists return false in the first param returned.
// If the category doesn't exist return false in the second param returned.
func (s *Service) CreateProp(prop *Prop) (bool, bool, error) {
	prop.AutoID()
	return s.crud.CreateWithRel(locService, s.cnt.StoreWithTimeout, prop)
}

// PutProp creates or updates a property into of the repository.
// If the property exists return true in the first param returned otherwise return false.
// If the category doesn't exist return false in the second param returned.
func (s *Service) PutProp(prop *Prop) (bool, bool, error) {
	return s.crud.PutWithRel(locService, s.cnt.StoreWithTimeout, prop)
}

// GetProp gets the property from storage
func (s *Service) GetProp(id xid.ID, prop *Prop) (bool, error) {
	ctx, cancel := s.cnt.StoreWithTimeout()
	ok, err := s.crud.Store().Get(ctx, entity.CatPropKey(id), prop)
	cancel()
	return ok, err
}

// LinkToProp links a category property with other category property or ente property
// of the child category. tPropID can be category property or ente property.
// The source an target must have the same type of data.
// The first parameter returned is true if the catPropID is found.
// The second parameter returned is true if the tPropID is found.
// The third parameter returned is true if the parent of tPropID is a child of the catPropID.
// The fourth parameter returned is true if the type of catPropID is equal to tPropID.
func (s *Service) LinkToProp(catPropID xid.ID, tPropID xid.ID) (bool, bool, bool, bool, relation.CatPropProp, error) {
	var scatProp Prop
	found, err := s.getProp(catPropID, &scatProp, "getting the source category property for linking")
	if err != nil {
		return false, false, false, false, relation.CatPropProp{}, err
	}
	if !found {
		return false, false, false, false, relation.CatPropProp{}, nil
	}

	found, tprop, err := s.propType(tPropID)
	if err != nil {
		return false, false, false, false, relation.CatPropProp{}, err
	}
	if !found {
		return true, false, false, false, relation.CatPropProp{}, nil
	}

	// Checks if the target property is child of the source property category
	ctx, cancel := s.cnt.StoreWithTimeout()
	found, err = s.crud.Store().Exists(ctx, storage.DLRKey(tprop.ParentKey(), scatProp.ParentKey()))
	cancel()
	if err != nil {
		return true, true, false, false, relation.CatPropProp{},
			s.cnt.Log.ErrWrap2(
				err,
				"checking if the target category property is child of the source property category",
				locService,
				logging.String("Source property", catPropID.String()),
				logging.String("Target property", tPropID.String()))
	}
	if !found {
		return true, true, false, false, relation.CatPropProp{}, nil
	}

	var txn storage.Txn

	// If the category property is not configured, this property is configured with the type
	// of the first property (category or ente)
	if scatProp.Type == entity.None {
		txn = storage.NewTxn(s.crud.Store())

		scatProp.Type = tprop.GetType()
		upd, err := s.crud.Store().Put(&scatProp)
		if err != nil {
			return true, true, true, false, relation.CatPropProp{},
				s.cnt.Log.ErrWrap2(
					err,
					"updating the type of category property before linking",
					locService,
					logging.String("Source property", catPropID.String()),
					logging.String("Target property", tPropID.String()))
		}
		txn.DoNotFound(upd)
	}

	if scatProp.Type != tprop.GetType() {
		return true, true, true, false, relation.CatPropProp{}, nil
	}

	// Link porperties and the same transaction updates the type of property
	cfound, pfound, link, err := s.crud.LinkTo(
		locService,
		s.cnt.StoreWithTimeout,
		txn,
		tprop.(storage.EntityRelation),
		entity.CatPropKey(catPropID),
		func(child storage.Entity) {
			switch p := child.(type) {
			case *Prop:
				p.catPropID = catPropID
			case *ente.Prop:
				p.CatPropID = catPropID
			}
		})
	if err != nil {
		return true, true, true, true, relation.CatPropProp{}, err
	}

	if !cfound || !pfound {
		return pfound, cfound, true, true, relation.CatPropProp{}, nil
	}

	return true, true, true, true, *link.(*relation.CatPropProp), nil
}

// propType gets the type of property (entity.TypeProp) and the parent identifier
func (s *Service) propType(tPropID xid.ID) (bool, entity.Property, error) {
	var prop entity.Property
	// I research the property type (category or ente)
	ctx, cancel := s.cnt.StoreWithTimeout()
	found, err := s.crud.Store().Exists(ctx, entity.CatPropKey(tPropID))
	cancel()
	if err != nil {
		return false, nil,
			s.cnt.Log.ErrWrap1(
				err,
				"it researching the the property type for linking",
				locService,
				logging.String("Property", tPropID.String()))
	}
	if found { // Category property
		var tcatProp Prop
		found, err := s.getProp(tPropID, &tcatProp, "getting the target category property for linking")
		if err != nil {
			return false, nil, err
		}
		if !found {
			return false, nil, err
		}
		prop = &tcatProp
	} else { // Ente property
		var tenteProp ente.Prop
		found, err := s.entesrv.GetProp(tPropID, &tenteProp)
		if err != nil {
			return false, nil,
				s.cnt.Log.ErrWrap1(
					err,
					"getting the target ente property for linking",
					locService,
					logging.String("Property", tPropID.String()))
		}
		if !found {
			return false, nil, nil
		}
		prop = &tenteProp
	}
	return true, prop, nil
}

func (s *Service) getProp(propID xid.ID, catsProp *Prop, errDesc string) (bool, error) {
	found, err := s.GetProp(propID, catsProp)
	if err != nil {
		return false, s.cnt.Log.ErrWrap1(
			err,
			errDesc,
			locService,
			logging.String("Property", propID.String()))
	}
	if !found {
		return false, nil
	}
	return true, nil
}
