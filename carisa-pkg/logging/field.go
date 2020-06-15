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

// logging values wrapper
package logging

// A FieldType indicates which member of the Field union struct should be used
type fieldType uint8

const (
	stringType fieldType = iota
	boolType
)

type Field struct {
	key     string
	tpy     fieldType
	boolV   bool
	stringV string
}

// String constructs a field that carries a string.
func String(key string, value string) Field {
	return Field{
		key:     key,
		tpy:     stringType,
		stringV: value,
	}
}

// Bool constructs a field that carries a bool.
func Bool(key string, value bool) Field {
	return Field{
		key:   key,
		tpy:   boolType,
		boolV: value,
	}
}
