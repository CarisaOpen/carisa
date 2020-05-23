/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and  limitations under the License.
 */

package encoding

import (
	"bytes"
	"encoding/gob"
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
