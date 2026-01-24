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
		cat Category
		ver Version
		fmt Format
	}
	testKeyHeaders := []struct {
		header   KeyHeader
		expected expected
	}{
		{
			header: NewKeyHeaderWith(Category(1), V1, Format(1)),
			expected: expected{
				cat: Category(1),
				ver: V1,
				fmt: Format(1),
			},
		},
		{
			header: NewKeyHeaderWith(Category(2), V1, Format(2)),
			expected: expected{
				cat: Category(2),
				ver: V1,
				fmt: Format(2),
			},
		},
		{
			header: NewKeyHeaderWith(Category(3), V1, Format(3)),
			expected: expected{
				cat: Category(3),
				ver: V1,
				fmt: Format(3),
			},
		},
		{
			header: NewKeyHeaderWith(Category(4), V1, Format(4)),
			expected: expected{
				cat: Category(4),
				ver: V1,
				fmt: Format(4),
			},
		},
	}
	for _, key := range testKeyHeaders {
		if key.header.Category() != key.expected.cat {
			t.Errorf("%v != %v", key.header.Category(), key.expected.cat)
		}
		if key.header.Version() != key.expected.ver {
			t.Errorf("%v != %v", key.header.Version(), key.expected.ver)
		}
		if key.expected.fmt != Format(0) {
			if key.header.Format() != key.expected.fmt {
				t.Errorf("%v != %v", key.header.Format(), key.expected.fmt)
			}
		}
	}
}
