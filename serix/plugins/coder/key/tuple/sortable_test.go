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
	"math"
	"testing"
)

// TestIntRoundTrip tests round-trip encoding/decoding of integers
func TestIntRoundTrip(t *testing.T) {
	testCases := []int64{
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

	for _, expected := range testCases {
		tpl := Tuple{expected}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Errorf("Pack failed for %d: %v", expected, err)
			continue
		}

		decoded, err := Unpack(encoded)
		if err != nil {
			t.Errorf("Unpack failed for %d: %v", expected, err)
			continue
		}

		if len(decoded) != 1 {
			t.Errorf("Expected 1 element, got %d", len(decoded))
			continue
		}

		actual, ok := decoded[0].(int64)
		if !ok {
			t.Errorf("Expected int64, got %T", decoded[0])
			continue
		}

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	}
}

// TestFloatRoundTrip tests round-trip encoding/decoding of floats
func TestFloatRoundTrip(t *testing.T) {
	testCases := []float64{
		math.Inf(-1),
		-1000.5,
		-1.0,
		-0.5,
		-0.0,
		0.0,
		0.5,
		1.0,
		1000.5,
		math.Inf(1),
	}

	for _, expected := range testCases {
		tpl := Tuple{expected}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Errorf("Pack failed for %f: %v", expected, err)
			continue
		}

		decoded, err := Unpack(encoded)
		if err != nil {
			t.Errorf("Unpack failed for %f: %v", expected, err)
			continue
		}

		if len(decoded) != 1 {
			t.Errorf("Expected 1 element, got %d", len(decoded))
			continue
		}

		actual, ok := decoded[0].(float64)
		if !ok {
			t.Errorf("Expected float64, got %T", decoded[0])
			continue
		}

		// Handle -0.0 vs 0.0 (both have same encoding)
		if math.Float64bits(actual) != math.Float64bits(expected) {
			t.Errorf("Expected %f (bits: %x), got %f (bits: %x)",
				expected, math.Float64bits(expected),
				actual, math.Float64bits(actual))
		}
	}
}

// TestFloatNaNRoundTrip tests that NaN roundtrips (bits may differ but should still be NaN)
func TestFloatNaNRoundTrip(t *testing.T) {
	tpl := Tuple{math.NaN()}
	encoded, err := tpl.Pack()
	if err != nil {
		t.Fatalf("Pack failed for NaN: %v", err)
	}

	decoded, err := Unpack(encoded)
	if err != nil {
		t.Fatalf("Unpack failed for NaN: %v", err)
	}

	if len(decoded) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(decoded))
	}

	actual, ok := decoded[0].(float64)
	if !ok {
		t.Fatalf("Expected float64, got %T", decoded[0])
	}

	if !math.IsNaN(actual) {
		t.Errorf("Expected NaN, got %f", actual)
	}
}

// TestStringRoundTrip tests round-trip encoding/decoding of strings
func TestStringRoundTrip(t *testing.T) {
	testCases := []string{
		"",
		"a",
		"aa",
		"b",
		"hello",
		"world",
		"a\x00b",           // embedded null
		"\x00",             // single null
		"\x00\x00",         // double null
		"test\x00test",     // null in middle
		"UTF-8: 日本語",    // UTF-8 characters
	}

	for _, expected := range testCases {
		tpl := Tuple{expected}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Errorf("Pack failed for %q: %v", expected, err)
			continue
		}

		decoded, err := Unpack(encoded)
		if err != nil {
			t.Errorf("Unpack failed for %q: %v", expected, err)
			continue
		}

		if len(decoded) != 1 {
			t.Errorf("Expected 1 element, got %d", len(decoded))
			continue
		}

		actual, ok := decoded[0].(string)
		if !ok {
			t.Errorf("Expected string, got %T", decoded[0])
			continue
		}

		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	}
}

// TestIntSortOrder tests that integer tuples sort in numeric order
func TestIntSortOrder(t *testing.T) {
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
		tpl := Tuple{v}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Fatalf("Pack failed for %d: %v", v, err)
		}
		encodings = append(encodings, encoded)
	}

	// Verify each pair is in correct order
	for i := 0; i < len(encodings)-1; i++ {
		cmp := bytes.Compare(encodings[i], encodings[i+1])
		if cmp >= 0 {
			t.Errorf("Sort order violation: %d (% x) should be < %d (% x), but bytes.Compare = %d",
				values[i], encodings[i], values[i+1], encodings[i+1], cmp)
		}
	}
}

// TestFloatSortOrder tests that float tuples sort in numeric order (excluding NaN)
func TestFloatSortOrder(t *testing.T) {
	// Values in expected sort order
	values := []float64{
		math.Inf(-1),
		-1000.5,
		-1.0,
		-0.5,
		-0.0, // Note: -0.0 and 0.0 encode the same
		0.0,
		0.5,
		1.0,
		1000.5,
		math.Inf(1),
	}

	// Encode all values
	var encodings [][]byte
	for _, v := range values {
		tpl := Tuple{v}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Fatalf("Pack failed for %f: %v", v, err)
		}
		encodings = append(encodings, encoded)
	}

	// Verify each pair is in correct order (or equal for -0.0 and 0.0)
	for i := 0; i < len(encodings)-1; i++ {
		cmp := bytes.Compare(encodings[i], encodings[i+1])
		if cmp > 0 {
			t.Errorf("Sort order violation: %f (% x) should be <= %f (% x), but bytes.Compare = %d",
				values[i], encodings[i], values[i+1], encodings[i+1], cmp)
		}
	}
}

// TestStringSortOrder tests that string tuples sort in lexicographic order
func TestStringSortOrder(t *testing.T) {
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
		tpl := Tuple{v}
		encoded, err := tpl.Pack()
		if err != nil {
			t.Fatalf("Pack failed for %q: %v", v, err)
		}
		encodings = append(encodings, encoded)
	}

	// Verify each pair is in correct order
	for i := 0; i < len(encodings)-1; i++ {
		cmp := bytes.Compare(encodings[i], encodings[i+1])
		if cmp >= 0 {
			t.Errorf("Sort order violation: %q (% x) should be < %q (% x), but bytes.Compare = %d",
				values[i], encodings[i], values[i+1], encodings[i+1], cmp)
		}
	}
}

// TestMixedTupleSortOrder tests that mixed-type tuples sort correctly
func TestMixedTupleSortOrder(t *testing.T) {
	// Test that tuples with different first elements sort by first element
	testCases := []struct {
		name   string
		tuples []Tuple
	}{
		{
			name: "int_then_string",
			tuples: []Tuple{
				{int64(1), "a"},
				{int64(1), "b"},
				{int64(2), "a"},
			},
		},
		{
			name: "string_then_int",
			tuples: []Tuple{
				{"a", int64(1)},
				{"a", int64(2)},
				{"b", int64(1)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode all tuples
			var encodings [][]byte
			for _, tpl := range tc.tuples {
				encoded, err := tpl.Pack()
				if err != nil {
					t.Fatalf("Pack failed for %v: %v", tpl, err)
				}
				encodings = append(encodings, encoded)
			}

			// Verify each pair is in correct order
			for i := 0; i < len(encodings)-1; i++ {
				cmp := bytes.Compare(encodings[i], encodings[i+1])
				if cmp >= 0 {
					t.Errorf("Sort order violation: %v (% x) should be < %v (% x), but bytes.Compare = %d",
						tc.tuples[i], encodings[i], tc.tuples[i+1], encodings[i+1], cmp)
				}
			}
		})
	}
}
