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
	}
	testKeyHeaders := []struct {
		header   KeyHeader
		expected expected
	}{
		{
			header: NewKeyHeaderWith(HeaderType(1), V1, ObjectType(1)),
			expected: expected{
				tp:  HeaderType(1),
				ver: V1,
				doc: ObjectType(1),
			},
		},
		{
			header: NewKeyHeaderWith(HeaderType(2), V1, ObjectType(2)),
			expected: expected{
				tp:  HeaderType(2),
				ver: V1,
				doc: ObjectType(2),
			},
		},
		{
			header: NewKeyHeaderWith(HeaderType(3), V1, ObjectType(3)),
			expected: expected{
				tp:  HeaderType(3),
				ver: V1,
				doc: ObjectType(3),
			},
		},
		{
			header: NewKeyHeaderWith(HeaderType(4), V1, ObjectType(4)),
			expected: expected{
				tp:  HeaderType(4),
				ver: V1,
				doc: ObjectType(4),
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
	}
}
