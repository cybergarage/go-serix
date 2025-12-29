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

package key

import (
	"bytes"
	"math"
	"testing"

	"github.com/cybergarage/go-serix/serix/document"
)

// SortableTest tests that the given coder produces sortable encodings for various key types.
func SortableTest(t *testing.T, coder document.KeyCoder) {
	t.Helper()

	t.Run("int", func(t *testing.T) {
		// Values in expected sort order
		values := []int64{
			math.MinInt64,
			-1000,
			-2,
			-1,
			0,
			1,
			2,
			1000,
			math.MaxInt64,
		}

		// Encode all values
		var encodings [][]byte
		for _, v := range values {
			key := document.NewKeyWith(v)
			encoded, err := coder.EncodeKey(key)
			if err != nil {
				t.Fatalf("Encode failed for %d: %v", v, err)
			}
			encodings = append(encodings, encoded)
		}

		// Verify each pair is in correct order
		for i := range len(encodings) - 1 {
			cmp := bytes.Compare(encodings[i], encodings[i+1])
			if cmp >= 0 {
				t.Errorf("Sort order violation: %d (% x) should be < %d (% x), but bytes.Compare = %d",
					values[i], encodings[i], values[i+1], encodings[i+1], cmp)
			}
		}
	})

	t.Run("float64", func(t *testing.T) {
		// Values in expected sort order
		values := []float64{
			math.Inf(-1),
			-1000.5,
			-1.0,
			-0.5,
			math.Copysign(0, -1), // Note: -0.0 and 0.0 encode the same
			0.0,
			0.5,
			1.0,
			1000.5,
			math.Inf(1),
		}

		// Encode all values
		var encodings [][]byte
		for _, v := range values {
			key := document.NewKeyWith(v)
			encoded, err := coder.EncodeKey(key)
			if err != nil {
				t.Fatalf("Encode failed for %f: %v", v, err)
			}
			encodings = append(encodings, encoded)
		}

		// Verify each pair is in correct order (or equal for -0.0 and 0.0)
		for i := range len(encodings) - 1 {
			cmp := bytes.Compare(encodings[i], encodings[i+1])
			if cmp > 0 {
				t.Errorf("Sort order violation: %f (% x) should be <= %f (% x), but bytes.Compare = %d",
					values[i], encodings[i], values[i+1], encodings[i+1], cmp)
			}
		}
	})

	t.Run("float64", func(t *testing.T) {
		// Values in expected sort order
		values := []string{
			"",
			"\x00",
			"\x00\x00",
			"a",
			"a\x00",
			"a\x00b",
			"aa",
			"aaa",
			"ab",
			"b",
			"ba",
			"hello",
			"world",
		}

		// Encode all values
		var encodings [][]byte
		for _, v := range values {
			key := document.NewKeyWith(v)
			encoded, err := coder.EncodeKey(key)
			if err != nil {
				t.Fatalf("Encode failed for %q: %v", v, err)
			}
			encodings = append(encodings, encoded)
		}

		// Verify each pair is in correct order
		for i := range len(encodings) - 1 {
			cmp := bytes.Compare(encodings[i], encodings[i+1])
			if cmp >= 0 {
				t.Errorf("Sort order violation: %q (% x) should be < %q (% x), but bytes.Compare = %d",
					values[i], encodings[i], values[i+1], encodings[i+1], cmp)
			}
		}
	})

	t.Run("float64", func(t *testing.T) {
		//  that tuples with different first elements sort by first element
		testCases := []struct {
			name string
			keys []document.Key
		}{
			{
				name: "int_then_string",
				keys: []document.Key{
					document.NewKeyWith(int64(1), "a"),
					document.NewKeyWith(int64(1), "b"),
					document.NewKeyWith(int64(2), "a"),
				},
			},
			{
				name: "string_then_int",
				keys: []document.Key{
					document.NewKeyWith("a", int64(1)),
					document.NewKeyWith("a", int64(2)),
					document.NewKeyWith("b", int64(1)),
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Encode all tuples
				var encodings [][]byte
				for _, key := range tc.keys {
					encoded, err := coder.EncodeKey(key)
					if err != nil {
						t.Fatalf("Encode failed for %v: %v", key, err)
					}
					encodings = append(encodings, encoded)
				}

				// Verify each pair is in correct order
				for i := range len(encodings) - 1 {
					cmp := bytes.Compare(encodings[i], encodings[i+1])
					if cmp >= 0 {
						t.Errorf("Sort order violation: %v (% x) should be < %v (% x), but bytes.Compare = %d",
							tc.keys[i], encodings[i], tc.keys[i+1], encodings[i+1], cmp)
					}
				}
			})
		}
	})
}
