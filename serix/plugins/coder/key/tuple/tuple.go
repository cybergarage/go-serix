// Copyright (C) 2022 The go-serix Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tuple

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-serix/serix/document"
)

// Tuple represents a tuple of elements that can be packed into a sortable byte sequence.
//
// COMPATIBILITY WARNING: The encoding format for int, float, and string types has been
// updated to support bytewise-sortable ordering for RocksDB. Keys encoded with the previous
// version are NOT compatible with this version and will not sort correctly.
type Tuple []any

const (
	markerNull   byte = 0x00
	markerTrue   byte = 0x01
	markerFalse  byte = 0x02
	markerInt    byte = 0x10
	markerUint   byte = 0x11
	markerFloat  byte = 0x20
	markerString byte = 0x30
	markerBytes  byte = 0x40
)

const (
	// String encoding uses escape sequences
	stringEscapeByte byte = 0x00
	stringEscapeNext byte = 0xFF
	stringTerminator byte = 0x00
	stringTermNext   byte = 0x00
)

// Pack encodes the tuple into a byte slice using gob encoding.
func (t Tuple) Pack() ([]byte, error) {
	packed, err := t.packSimple()
	if err != nil {
		return nil, err
	}
	return packed, nil
}

func (t Tuple) packSimple() ([]byte, error) {
	var buf bytes.Buffer

	for _, elem := range t {
		switch v := elem.(type) {
		case nil:
			buf.WriteByte(markerNull)
		case bool:
			if v {
				buf.WriteByte(markerTrue)
			} else {
				buf.WriteByte(markerFalse)
			}
		case int, int8, int16, int32, int64:
			buf.WriteByte(markerInt)
			var tv int64
			if err := safecast.ToInt64(v, &tv); err != nil {
				return nil, err
			}
			// Sortable int encoding: flip sign bit to make negative values sort before positive
			sortable := uint64(tv) ^ (1 << 63)
			binary.Write(&buf, binary.BigEndian, sortable)
		case uint, uint8, uint16, uint32, uint64:
			buf.WriteByte(markerUint)
			var tv uint64
			if err := safecast.ToUint64(v, &tv); err != nil {
				return nil, err
			}
			binary.Write(&buf, binary.BigEndian, tv)
		case float32, float64:
			buf.WriteByte(markerFloat)
			var tv float64
			if err := safecast.ToFloat64(v, &tv); err != nil {
				return nil, err
			}
			// Sortable float encoding: transform IEEE754 bits to be bytewise sortable
			bits := math.Float64bits(tv)
			if bits&(1<<63) != 0 {
				// Negative: flip all bits
				bits = ^bits
			} else {
				// Positive: flip only sign bit
				bits = bits ^ (1 << 63)
			}
			binary.Write(&buf, binary.BigEndian, bits)
		case string:
			buf.WriteByte(markerString)
			data := []byte(v)
			// Sortable string encoding: escape 0x00 bytes and terminate with 0x00 0x00
			for _, b := range data {
				if b == stringEscapeByte {
					buf.WriteByte(stringEscapeByte)
					buf.WriteByte(stringEscapeNext)
				} else {
					buf.WriteByte(b)
				}
			}
			// Terminator
			buf.WriteByte(stringTerminator)
			buf.WriteByte(stringTermNext)
		case []byte:
			buf.WriteByte(markerBytes)
			binary.Write(&buf, binary.BigEndian, uint32(len(v)))
			buf.Write(v)
		default:
			// Convert unknown types to strings
			str := fmt.Sprintf("%v", v)
			buf.WriteByte(markerString)
			data := []byte(str)
			// Sortable string encoding: escape 0x00 bytes and terminate with 0x00 0x00
			for _, b := range data {
				if b == stringEscapeByte {
					buf.WriteByte(stringEscapeByte)
					buf.WriteByte(stringEscapeNext)
				} else {
					buf.WriteByte(b)
				}
			}
			// Terminator
			buf.WriteByte(stringTerminator)
			buf.WriteByte(stringTermNext)
		}
	}
	return buf.Bytes(), nil
}

// Unpack decodes a byte slice into a tuple.
func Unpack(data []byte) (Tuple, error) {
	return unpackSimple(data)
}

func unpackSimple(data []byte) (Tuple, error) {
	buf := bytes.NewBuffer(data)
	var tuple Tuple

	var err error
	for buf.Len() > 0 {
		var marker byte
		marker, err = buf.ReadByte()
		if err != nil {
			break
		}

		switch marker {
		case markerNull:
			tuple = append(tuple, nil)
		case markerTrue:
			tuple = append(tuple, true)
		case markerFalse:
			tuple = append(tuple, false)
		case markerInt:
			var sortable uint64
			err := binary.Read(buf, binary.BigEndian, &sortable)
			if err != nil {
				return nil, err
			}
			// Reverse sortable int encoding: flip sign bit back
			val := int64(sortable ^ (1 << 63))
			tuple = append(tuple, val)
		case markerUint:
			var val uint64
			err := binary.Read(buf, binary.BigEndian, &val)
			if err != nil {
				return nil, err
			}
			tuple = append(tuple, val)
		case markerFloat:
			var sortableBits uint64
			err := binary.Read(buf, binary.BigEndian, &sortableBits)
			if err != nil {
				return nil, err
			}
			// Reverse sortable float encoding
			var bits uint64
			if sortableBits&(1<<63) != 0 {
				// Was positive: flip only sign bit
				bits = sortableBits ^ (1 << 63)
			} else {
				// Was negative: flip all bits
				bits = ^sortableBits
			}
			val := math.Float64frombits(bits)
			tuple = append(tuple, val)
		case markerString:
			// Sortable string decoding: read until 0x00 0x00 terminator, unescape 0x00 0xFF -> 0x00
			var data []byte
			for {
				b, err := buf.ReadByte()
				if err != nil {
					return nil, fmt.Errorf("unexpected end while reading string")
				}
				if b == stringEscapeByte {
					next, err := buf.ReadByte()
					if err != nil {
						return nil, fmt.Errorf("unexpected end after escape byte")
					}
					if next == stringEscapeNext {
						// Escaped 0x00
						data = append(data, stringEscapeByte)
					} else if next == stringTermNext {
						// Found terminator 0x00 0x00
						break
					} else {
						return nil, fmt.Errorf("invalid escape sequence: 0x00 0x%02x", next)
					}
				} else {
					data = append(data, b)
				}
			}
			tuple = append(tuple, string(data))
		case markerBytes:
			var length uint32
			err := binary.Read(buf, binary.BigEndian, &length)
			if err != nil {
				return nil, err
			}
			data := make([]byte, length)
			_, err = buf.Read(data)
			if err != nil {
				return nil, err
			}
			tuple = append(tuple, data)
		default:
			return nil, fmt.Errorf("unknown marker: %02x", marker)
		}
	}

	return tuple, err
}

func newTupleWith(key document.Key) (Tuple, error) {
	tpl := make(Tuple, len(key))
	copy(tpl, key)
	return tpl, nil
}

func newKeyWith(tpl Tuple) document.Key {
	key := make(document.Key, len(tpl))
	copy(key, tpl)
	return key
}
