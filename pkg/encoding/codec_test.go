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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	S string
	I int
	B bool
}

func TestEncoding_EncodeDecode(t *testing.T) {
	data := T{
		S: "String",
		I: 1,
		B: true,
	}

	encode, err := Encode(data)
	if assert.NoError(t, err, "unexpected encode error") {
		var decode T
		err = Decode(encode, &decode)
		if assert.NoError(t, err, "unexpected decode error") {
			assert.Equal(t, data, decode, "Encode and decode values are not equal")
		}
	}
}

func TestEncoding_DecodeByte(t *testing.T) {
	data := T{
		S: "String",
		I: 1,
		B: true,
	}

	encode, err := Encode(data)
	assert.NoError(t, err, "unexpected encode error")

	var decode T
	err = DecodeByte([]byte(encode), &decode)
	assert.NoError(t, err, "unexpected decode error")

	assert.Equal(t, data, decode, "Encode and decode values are not equal")
}

func TestEncoding_Encodeuint32Desc(t *testing.T) {
	codec, _ := EncodeUI32Desc(1234567890)
	assert.Equal(t, "ihgfedcbaj", codec)
}

func TestEncoding_EncodeStrDesc(t *testing.T) {
	codec, _ := EncodeStrDesc("0123456789")
	assert.Equal(t, "jihgfedcba", codec)
}

func TestEncoding_EncodeStrDesc_Err(t *testing.T) {
	_, err := EncodeStrDesc("1234A5789")
	assert.Error(t, err)
}

func TestEncoding_EncodeStrDesc_Sort(t *testing.T) {
	array := make([]string, 4)
	array[0], _ = EncodeStrDesc("3000256")
	array[1], _ = EncodeStrDesc("4000266")
	array[2], _ = EncodeStrDesc("4000256")
	array[3], _ = EncodeStrDesc("4010256")
	sort.Strings(array)
	assert.Equal(t, []string{"3232", "23232"}, array)
}
