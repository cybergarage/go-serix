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
	"encoding/gob"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
)

// Tuple represents a tuple of elements that can be packed into a sortable byte sequence
type Tuple []interface{}

// Pack encodes the tuple into a byte slice using gob encoding
func (t Tuple) Pack() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	// Encode each element with type information
	elements := make([]encodedElement, len(t))
	for i, elem := range t {
		elements[i] = encodeElementWithType(elem)
	}

	err := encoder.Encode(elements)
	if err != nil {
		// Fallback to simple binary encoding for basic types
		return t.packSimple()
	}
	return buf.Bytes()
}

type encodedElement struct {
	Type  string
	Value interface{}
}

func encodeElementWithType(elem interface{}) encodedElement {
	switch v := elem.(type) {
	case nil:
		return encodedElement{Type: "nil", Value: nil}
	case bool:
		return encodedElement{Type: "bool", Value: v}
	case int:
		return encodedElement{Type: "int", Value: int64(v)}
	case int8:
		return encodedElement{Type: "int8", Value: int64(v)}
	case int16:
		return encodedElement{Type: "int16", Value: int64(v)}
	case int32:
		return encodedElement{Type: "int32", Value: int64(v)}
	case int64:
		return encodedElement{Type: "int64", Value: v}
	case uint:
		return encodedElement{Type: "uint", Value: uint64(v)}
	case uint8:
		return encodedElement{Type: "uint8", Value: uint64(v)}
	case uint16:
		return encodedElement{Type: "uint16", Value: uint64(v)}
	case uint32:
		return encodedElement{Type: "uint32", Value: uint64(v)}
	case uint64:
		return encodedElement{Type: "uint64", Value: v}
	case float32:
		return encodedElement{Type: "float32", Value: float64(v)}
	case float64:
		return encodedElement{Type: "float64", Value: v}
	case string:
		return encodedElement{Type: "string", Value: v}
	case []byte:
		return encodedElement{Type: "bytes", Value: v}
	default:
		return encodedElement{Type: "string", Value: fmt.Sprintf("%v", v)}
	}
}

func (t Tuple) packSimple() []byte {
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
			val := toInt64(v)
			binary.Write(&buf, binary.BigEndian, val)
		case uint, uint8, uint16, uint32, uint64:
			buf.WriteByte(0x11) // uint marker
			val := toUint64(v)
			binary.Write(&buf, binary.BigEndian, val)
		case float32, float64:
			buf.WriteByte(0x20) // float marker
			val := toFloat64(v)
			binary.Write(&buf, binary.BigEndian, val)
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
	return buf.Bytes()
}

// Unpack decodes a byte slice into a tuple
func Unpack(data []byte) (Tuple, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	var elements []encodedElement
	err := decoder.Decode(&elements)
	if err != nil {
		// Fallback to simple decoding
		return unpackSimple(data)
	}

	tuple := make(Tuple, len(elements))
	for i, elem := range elements {
		tuple[i] = decodeElementWithType(elem)
	}

	return tuple, nil
}

func decodeElementWithType(elem encodedElement) interface{} {
	switch elem.Type {
	case "nil":
		return nil
	case "bool":
		return elem.Value.(bool)
	case "int":
		return int(elem.Value.(int64))
	case "int8":
		return int8(elem.Value.(int64))
	case "int16":
		return int16(elem.Value.(int64))
	case "int32":
		return int32(elem.Value.(int64))
	case "int64":
		return elem.Value.(int64)
	case "uint":
		return uint(elem.Value.(uint64))
	case "uint8":
		return uint8(elem.Value.(uint64))
	case "uint16":
		return uint16(elem.Value.(uint64))
	case "uint32":
		return uint32(elem.Value.(uint64))
	case "uint64":
		return elem.Value.(uint64)
	case "float32":
		return float32(elem.Value.(float64))
	case "float64":
		return elem.Value.(float64)
	case "string":
		return elem.Value.(string)
	case "bytes":
		return elem.Value.([]byte)
	default:
		return elem.Value
	}
}

func unpackSimple(data []byte) (Tuple, error) {
	buf := bytes.NewBuffer(data)
	var tuple Tuple

	for buf.Len() > 0 {
		marker, err := buf.ReadByte()
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

	return tuple, nil
}

func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	default:
		return 0
	}
}

func toUint64(v interface{}) uint64 {
	switch val := v.(type) {
	case uint:
		return uint64(val)
	case uint8:
		return uint64(val)
	case uint16:
		return uint64(val)
	case uint32:
		return uint64(val)
	case uint64:
		return val
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float32:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
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
