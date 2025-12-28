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

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-serix/serix/document"
)

// Tuple represents a tuple of elements that can be packed into a sortable byte sequence.
type Tuple []any

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
			buf.WriteByte(0x00) // null marker
		case bool:
			if v {
				buf.WriteByte(0x01) // true marker
			} else {
				buf.WriteByte(0x02) // false marker
			}
		case int, int8, int16, int32, int64:
			buf.WriteByte(0x10) // int marker
			var tv int64
			if err := safecast.ToInt64(v, &tv); err != nil {
				return nil, err
			}
			binary.Write(&buf, binary.BigEndian, tv)
		case uint, uint8, uint16, uint32, uint64:
			buf.WriteByte(0x11) // uint marker
			var tv uint64
			if err := safecast.ToUint64(v, &tv); err != nil {
				return nil, err
			}
			binary.Write(&buf, binary.BigEndian, tv)
		case float32, float64:
			buf.WriteByte(0x20) // float marker
			var tv float64
			if err := safecast.ToFloat64(v, &tv); err != nil {
				return nil, err
			}
			binary.Write(&buf, binary.BigEndian, tv)
		case string:
			buf.WriteByte(0x30) // string marker
			data := []byte(v)
			binary.Write(&buf, binary.BigEndian, uint32(len(data)))
			buf.Write(data)
		case []byte:
			buf.WriteByte(0x40) // bytes marker
			binary.Write(&buf, binary.BigEndian, uint32(len(v)))
			buf.Write(v)
		default:
			// Convert unknown types to strings
			str := fmt.Sprintf("%v", v)
			buf.WriteByte(0x30) // string marker
			data := []byte(str)
			binary.Write(&buf, binary.BigEndian, uint32(len(data)))
			buf.Write(data)
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
		case 0x00: // null
			tuple = append(tuple, nil)
		case 0x01: // true
			tuple = append(tuple, true)
		case 0x02: // false
			tuple = append(tuple, false)
		case 0x10: // int
			var val int64
			err := binary.Read(buf, binary.BigEndian, &val)
			if err != nil {
				return nil, err
			}
			tuple = append(tuple, val)
		case 0x11: // uint
			var val uint64
			err := binary.Read(buf, binary.BigEndian, &val)
			if err != nil {
				return nil, err
			}
			tuple = append(tuple, val)
		case 0x20: // float
			var val float64
			err := binary.Read(buf, binary.BigEndian, &val)
			if err != nil {
				return nil, err
			}
			tuple = append(tuple, val)
		case 0x30: // string
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
			tuple = append(tuple, string(data))
		case 0x40: // bytes
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
