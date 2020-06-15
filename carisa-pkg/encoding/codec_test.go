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
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	S string
	I int
	B bool
}

func TestEncodeDecode(t *testing.T) {
	data := T{
		S: "String",
		I: 1,
		B: true,
	}

	encode, err := Encode(data)
	assert.NoError(t, err, "unexpected encode error")

	var decode T
	err = Decode(encode, &decode)
	assert.NoError(t, err, "unexpected decode error")

	assert.Equal(t, data, decode, "Encode and decode values are not equal")
}
