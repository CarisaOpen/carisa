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

const MaxUint8 = ^uint8(0)
const MaxUint32 = ^uint32(0)

func TestSimpleCodec_EncodeDecode(t *testing.T) {
	encode := SimpleEncoder{}
	encode.WriteUint32(uint32(34))
	encode.WriteUint8(uint8(35))
	encode.WriteBytes([]byte("stringsssssss"))
	encode.WriteUint32(MaxUint32)
	encode.WriteUint8(MaxUint8)
	encode.WriteBytes([]byte("strings"))
	encode.WriteUint32(uint32(3141241231))
	encode.WriteUint8(uint8(120))
	encode.WriteBytes([]byte("str"))

	buffer := encode.Bytes()
	decode := NewSimpleDecoder(buffer)
	res32, _ := decode.ReadUint32()
	assert.Equal(t, uint32(34), res32)
	res8, _ := decode.ReadUint8()
	assert.Equal(t, uint8(35), res8)
	resb, _ := decode.ReadBytes()
	assert.Equal(t, "stringsssssss", string(resb))
	res32, _ = decode.ReadUint32()
	assert.Equal(t, MaxUint32, res32)
	res8, _ = decode.ReadUint8()
	assert.Equal(t, MaxUint8, res8)
	resb, _ = decode.ReadBytes()
	assert.Equal(t, "strings", string(resb))
	res32, _ = decode.ReadUint32()
	assert.Equal(t, uint32(3141241231), res32)
	res8, _ = decode.ReadUint8()
	assert.Equal(t, uint8(120), res8)
	resb, _ = decode.ReadBytes()
	assert.Equal(t, string("str"), string(resb))
}
