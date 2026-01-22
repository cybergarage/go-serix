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

package kv

import (
	"testing"
)

func TestKeyHeader(t *testing.T) {
	type expected struct {
		tp  HeaderType
		ver Version
		doc ObjectType
		idx DocumentSubType
	}
	testKeyHeaders := []struct {
		header   KeyHeader
		expected expected
	}{
		{
			header: DatabaseKeyHeader,
			expected: expected{
				tp:  ObjectType(1),
				ver: V1,
				doc: CBOR,
				idx: DocumentSubType(0),
			},
		},
		{
			header: CollectionKeyHeader,
			expected: expected{
				tp:  ObjectType(2),
				ver: V1,
				doc: CBOR,
				idx: DocumentSubType(0),
			},
		},
		{
			header: DocumentKeyHeader,
			expected: expected{
				tp:  ObjectType(3),
				ver: V1,
				doc: CBOR,
				idx: DocumentSubType(0),
			},
		},
		{
			header: IndexKeyHeader,
			expected: expected{
				tp:  ObjectType(4),
				ver: V1,
				doc: ObjectType(0),
				idx: DocumentSubType(2),
			},
		},
	}
	for _, key := range testKeyHeaders {
		if key.header.Type() != key.expected.tp {
			t.Errorf("%v != %v", key.header.Type(), key.expected.tp)
		}
		if key.header.Version() != key.expected.ver {
			t.Errorf("%v != %v", key.header.Version(), key.expected.ver)
		}
		if key.expected.doc != ObjectType(0) {
			if key.header.ObjectType() != key.expected.doc {
				t.Errorf("%v != %v", key.header.ObjectType(), key.expected.doc)
			}
		}
		if key.expected.idx != DocumentSubType(0) {
			if key.header.DocumentSubType() != key.expected.idx {
				t.Errorf("%v != %v", key.header.DocumentSubType(), key.expected.idx)
			}
		}
	}
}
