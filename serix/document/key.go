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

import (
	"fmt"

	"github.com/cybergarage/go-safecast/safecast"
)

// Key represents an unique key for a document object.
type Key []any

// NewKey returns a new blank key.
func NewKey() Key {
	return Key{}
}

// NewKeyWith returns a new key from the specified key elements.
func NewKeyWith(elems ...any) Key {
	elemArray := make([]any, len(elems))
	copy(elemArray, elems)
	return elemArray
}

// Elements returns all elements of the key.
func (key Key) Elements() []any {
	return key
}

// Len returns the number of elements of the key.
func (key Key) Len() int {
	return len(key)
}

// Equal reports whether the two keys are equal.
func (key Key) Equal(other Key) bool {
	for n, elem := range key {
		if !safecast.Equal(elem, other[n]) {
			return false
		}
	}
	return true
}

// String returns a string representation of the key.
func (key Key) String() string {
	var s string
	for _, elem := range key {
		s += fmt.Sprintf("%v", elem)
	}
	return s
}
