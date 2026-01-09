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
	"github.com/cybergarage/go-serix/serixtest/document/object"
)

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
	/*
		cascadeBinaryCoderTest := func(t *testing.T, compressor document.ObjectCoder) {
			t.Helper()

			if compressor.Type() != document.ObjectCompressor {
				t.Fatalf("expected document.ObjectCompressor, got %v", compressor.Type())
			}

			serializer := []document.ObjectCoder{}

			mgr := plugins.NewManager()
			for _, plugin := range mgr.ObjectCoders() {
				if compressor.Type() == document.ObjectSerializer {
					serializer = append(serializer, plugin)
				}
			}

			for _, serializer := range serializer {
				t.Run(serializer.Name()+" + "+compressor.Name(), func(t *testing.T) {
					coder := document.NewChainCorder(serializer, compressor)
					ObjectCompressorSuite(t, coder)
				})
			}

		}
	*/
	testFuncs := []struct {
		name string
		fn   func(*testing.T, document.ObjectCoder)
	}{
		{"single", object.BinaryCoderTest},
		// {"cascade", cascadeBinaryCoderTest},
	}

	for _, testFunc := range testFuncs {
		t.Run(testFunc.name, func(t *testing.T) {
			testFunc.fn(t, coder)
		})
	}
}
