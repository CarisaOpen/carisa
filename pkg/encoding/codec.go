/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 */

package encoding

import (
	"bytes"
	"encoding/gob"

	"github.com/carisa/pkg/strings"

	"github.com/pkg/errors"
)

// Encode encodes data to string using encoding gob
func Encode(data interface{}) (string, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Decode decodes data from string using decoding gob
func Decode(encode string, data interface{}) error {
	enc := gob.NewDecoder(bytes.NewBufferString(encode))
	return enc.Decode(data)
}

// Decode decodes data from bytes using decoding gob
func DecodeByte(encode []byte, data interface{}) error {
	enc := gob.NewDecoder(bytes.NewBuffer(encode))
	return enc.Decode(data)
}

// EncodeUI32Desc encodes a uint32 into letters
// that allow the uint32 to be sorted in descending order
// by applying ascending order
func EncodeUI32Desc(u uint32) (string, error) {
	return EncodeStrDesc(strings.Convertuint32(u))
}

// EncodeStrDesc encodes a string with positive real numbers into letters
// that allow the source number to be sorted in descending order
// by applying ascending order
func EncodeStrDesc(str string) (string, error) {
	const ascii0 = 48
	const ascii9 = 57
	const asciij = 106

	b := []byte(str)
	for i, c := range b {
		if !(c >= ascii0 && c <= ascii9) {
			return "", errors.New("The parameter must be a positive real number")
		}
		b[i] = asciij - (c - ascii0)
	}
	return strings.ConvertBytes(b), nil
}
