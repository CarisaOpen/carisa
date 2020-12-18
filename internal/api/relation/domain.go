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

package relation

import "github.com/carisa/pkg/strings"

// Link name. Must be unique
const (
	InstSpaceLn   = "IS"
	SpaceCatLn    = "SC"
	EntePropLn    = "EP"
	CatCatLn      = "CC"
	CatEnteLn     = "CE"
	CatPropLn     = "CP"
	SpaceEnteLn   = "SE"
	CatPropPropLn = "CPP"
)

// InstSpace represents the link between instance.Instance and space.Space
type InstSpace struct {
	ID      string `json:"-"`
	Name    string `json:"name"`
	SpaceID string `json:"spaceId"`
}

func (l *InstSpace) ToString() string {
	return strings.Concat("inst-space-link: ID:", l.Key(), ", Name:", l.Name)
}

func (l *InstSpace) Key() string {
	return l.ID
}

// SpaceEnte represents the link between space.Space and ente.Ente
type SpaceEnte struct {
	ID     string `json:"-"`
	Name   string `json:"name"`
	EnteID string `json:"enteId"`
}

func (s *SpaceEnte) ToString() string {
	return strings.Concat("space-ente-link: ID:", s.Key(), ", Name:", s.Name)
}

func (s *SpaceEnte) Key() string {
	return s.ID
}

// EnteProp represents the link between ente.Ente and her properties
type EnteProp struct {
	ID         string `json:"-"`
	Name       string `json:"name"`
	EntePropID string `json:"entePropId"`
}

func (s *EnteProp) ToString() string {
	return strings.Concat("ente-enteprop-link: ID:", s.Key(), ", Name:", s.Name)
}

func (s *EnteProp) Key() string {
	return s.ID
}

// SpaceCategory represents the link between space.Space and category.Category
type SpaceCategory struct {
	ID    string `json:"-"`
	Name  string `json:"name"`
	CatID string `json:"categoryId"`
}

func (s *SpaceCategory) ToString() string {
	return strings.Concat("space-category-link: ID:", s.Key(), ", Name:", s.Name)
}

func (s *SpaceCategory) Key() string {
	return s.ID
}

// Hierarchy represents the link between category.Category and others category.Category or ente.Ente
type Hierarchy struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	LinkID   string `json:"linkId"`
	Category bool   `json:"category"` // Category=false the hierarchy link to a ente.Ente
}

func (h *Hierarchy) ToString() string {
	return strings.Concat("hierarchy-link: ID:", h.Key(), ", Name:", h.Name)
}

func (h *Hierarchy) Key() string {
	return h.ID
}

// CategoryProp represents the link between category.Category and her properties
type CategoryProp struct {
	ID        string `json:"-"`
	Name      string `json:"name"`
	CatPropID string `json:"categoryPropId"`
}

func (c *CategoryProp) ToString() string {
	return strings.Concat("category-catprop-link: ID:", c.Key(), ", Name:", c.Name)
}

func (c *CategoryProp) Key() string {
	return c.ID
}

// CatPropProp represents the link between category.Category property and others category property or ente.Ente property
// of the children
type CatPropProp struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	PropID   string `json:"propertyId"`
	Category bool   `json:"category"` // Category=false the property links to a ente.Ente property
}

func (l *CatPropProp) ToString() string {
	return strings.Concat("catprop-link: ID:", l.Key(), ", Name:", l.Name)
}

func (l *CatPropProp) Key() string {
	return l.ID
}

// PlatformPlugin represents the link between platform and plugin.Prototype
type PlatformPlugin struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	ProtoID  string `json:"protoId"`
	Category string `json:"category"` // Category is the type of plugin (query, etc)
}

func (p *PlatformPlugin) ToString() string {
	return strings.Concat("platform-plugin-link: ID:", p.Key(), ", Name:", p.Name)
}

func (p *PlatformPlugin) Key() string {
	return p.ID
}

// PlatformInstance represents the link between platform and plugin instance.Instance
type PlatformInstance struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	InstID   string `json:"instanceId"`
	Category string `json:"category"` // Category is the type of plugin (query, etc)
}

func (p *PlatformInstance) ToString() string {
	return strings.Concat("platform-instance-link: ID:", p.Key(), ", Name:", p.Name)
}

func (p *PlatformInstance) Key() string {
	return p.ID
}
