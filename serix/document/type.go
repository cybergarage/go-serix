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

package document

// CoderType represents the type of a coder.
type CoderType int

const (
	// KeySerializer represents a key serializer coder type.
	KeySerializer CoderType = 0x00
	// ObjectSerializer represents an object serializer coder type.
	ObjectSerializer CoderType = 0x01
	// ObjectCompressor represents an object compressor coder type.
	ObjectCompressor CoderType = 0x02
)

// String returns the string representation of the coder type.
func (ct CoderType) String() string {
	switch ct {
	case KeySerializer:
		return "key_serializer"
	case ObjectSerializer:
		return "object_serializer"
	case ObjectCompressor:
		return "object_compressor"
	default:
		return "unknown"
	}
}
