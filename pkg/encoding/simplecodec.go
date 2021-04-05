/*
 *  Copyright 2019-2022 the original author or authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package encoding

import (
	"bytes"
	"io"
)

// SimpleEncoder represents a binary encoder to primitive structs (no escape to heap),
// Based on https://github.com/kelindar/binary/blob/master/encoder.go
type SimpleEncoder struct {
	scratch [10]byte
	out     bytes.Buffer
	err     error
}

// write writes the contents of p into the buffer.
func (e *SimpleEncoder) write(v []byte) {
	if e.err == nil {
		_, e.err = e.out.Write(v)
	}
}

// write writes a set of bytes into the buffer. Use for variable size
func (e *SimpleEncoder) WriteBytes(v []byte) {
	if e.err == nil {
		e.WriteUint32(uint32(len(v)))
		e.write(v)
	}
}

// WriteUint8 writes to Uint8
func (e *SimpleEncoder) WriteUint8(v uint8) {
	e.scratch[0] = v
	e.write(e.scratch[:1])
}

// WriteUint32 writes to Uint32
func (e *SimpleEncoder) WriteUint32(v uint32) {
	e.scratch[0] = byte(v)
	e.scratch[1] = byte(v >> 8)
	e.scratch[2] = byte(v >> 16)
	e.scratch[3] = byte(v >> 24)
	e.write(e.scratch[:4])
}

// String writes to Uint8
func (e *SimpleEncoder) String() string {
	return e.out.String()
}

// Buffer writes to Uint8
func (e *SimpleEncoder) Bytes() []byte {
	return e.out.Bytes()
}

type SimpleDecoder struct {
	s reader // Not using the interface for better inlining
}

// NewSimpleDecoder creates a binary decoder (no escape to heap).
func NewSimpleDecoder(buffer []byte) SimpleDecoder {
	return SimpleDecoder{
		s: reader{buffer, 0},
	}
}

// ReadBytes reads a set of bytes
func (d *SimpleDecoder) ReadBytes() (out []byte, err error) {
	len, err := d.ReadUint32()
	if err != nil {
		return nil, err
	}
	return d.s.slice(len)
}

// ReadUint32 reads a uint32
func (d *SimpleDecoder) ReadUint8() (out uint8, err error) {
	var b []byte
	if b, err = d.s.slice(1); err == nil {
		_ = b[0] // bounds check hint to compiler
		out = b[0]
	}
	return
}

// ReadUint32 reads a uint32
func (d *SimpleDecoder) ReadUint32() (out uint32, err error) {
	var b []byte
	if b, err = d.s.slice(4); err == nil {
		_ = b[3] // bounds check hint to compiler
		out = uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	}
	return
}

// A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner interfaces by reading from
// a byte slice.
// Unlike a Buffer, a Reader is read-only and supports seeking.
type reader struct {
	s []byte
	i int64 // current reading index
}

// slice selects a sub-slice of next bytes. This is similar to Read() but does not
// actually perform a copy, but simply uses the underlying slice (if available) and
// returns a sub-slice pointing to the same array. Since this requires access
// to the underlying data, this is only available for our default reader.
func (r *reader) slice(n uint32) ([]byte, error) {
	if r.i+int64(n) > int64(len(r.s)) {
		return nil, io.EOF
	}

	cur := r.i
	r.i += int64(n)
	return r.s[cur:r.i], nil
}
