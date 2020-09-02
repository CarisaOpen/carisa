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

// Link represents the link between instance and space
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

// Link represents the link between space and ente
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

// Link represents the link between ente and properties
type EnteEnteProp struct {
	ID         string `json:"-"`
	Name       string `json:"name"`
	EntePropID string `json:"entePropId"`
}

func (s *EnteEnteProp) ToString() string {
	return strings.Concat("ente-enteprop-link: ID:", s.Key(), ", Name:", s.Name)
}

func (s *EnteEnteProp) Key() string {
	return s.ID
}
