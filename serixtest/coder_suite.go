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

package serixtest

import (
	_ "embed"
	"testing"

	"github.com/cybergarage/go-serix/serix/document"
	"github.com/cybergarage/go-serix/serixtest/document/key"
	"github.com/cybergarage/go-serix/serixtest/document/object"
)

// KeyCoderSuite tests the encoding and decoding of keys using the provided KeyCoder.
func KeyCoderSuite(t *testing.T, coder document.KeyCoder) {
	t.Helper()

	tests := []struct {
		name string
		test func(t *testing.T, coder document.KeyCoder)
	}{
		{
			name: "RoundTripKeyTest",
			test: key.RoundTripKeyTest,
		},
		{
			name: "SortableKeyTest",
			test: key.SortableKeyTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, coder)
		})
	}
}

// ObjectSerializerSuite tests the specified document coder.
func ObjectSerializerSuite(t *testing.T, coder document.ObjectCoder) {
	t.Helper()

	testFuncs := []struct {
		name string
		fn   func(*testing.T, document.ObjectCoder)
	}{
		{"primitive", object.PrimitiveObjectTest},
		{"array", object.ArrayObjectTest},
		{"map", object.MapObjectTest},
	}

	for _, testFunc := range testFuncs {
		t.Run(testFunc.name, func(t *testing.T) {
			testFunc.fn(t, coder)
		})
	}
}

// ObjectCompressorSuite tests the binary encoding and decoding using the provided Coder.
func ObjectCompressorSuite(t *testing.T, coder document.ObjectCoder) {
	t.Helper()

	// cascadeCompressorSuite := func(t *testing.T, coder document.ObjectCoder) {
	// 	t.Helper()
	// 	plugins.CascadeCompressorSuite(t, coder, ObjectSerializerSuite)
	// }

	testFuncs := []struct {
		name string
		fn   func(*testing.T, document.ObjectCoder)
	}{
		{"single", object.BinaryCoderTest},
		// {"cascade", cascadeCompressorSuite},
	}

	for _, testFunc := range testFuncs {
		t.Run(testFunc.name, func(t *testing.T) {
			testFunc.fn(t, coder)
		})
	}
}
