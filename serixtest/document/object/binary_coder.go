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

package object

import (
	"bytes"
	"testing"

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-serix/serix/document"
)

// BinaryCoderTest tests the specified binary coder.
func BinaryCoderTest(t *testing.T, coder document.ObjectCoder) {
	t.Helper()

	tests := []struct {
		name string
		obj  any
	}{
		{
			name: "empty byte data",
			obj:  []byte{},
		},
		{
			name: "empty string data",
			obj:  "",
		},
		{
			name: "data with byte content",
			obj:  []byte("This is test content"),
		},
		{
			name: "data with string content",
			obj:  "This is test content",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			err := coder.EncodeObject(w, test.obj)
			if err != nil {
				t.Fatalf("failed to encode: %v", err)
			}

			r := bytes.NewBuffer(w.Bytes())
			decordedObj, err := coder.DecodeObject(r)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}

			switch testObj := test.obj.(type) {
			case []byte:
				var decordedObjBytes []byte
				if err := safecast.ToBytes(decordedObj, &decordedObjBytes); err != nil {
					t.Fatalf("failed to cast decoded object to []byte: %v", err)
				}
				if !bytes.Equal(testObj, decordedObjBytes) {
					t.Errorf("decoded data does not match original: expected %v, got %v", testObj, decordedObjBytes)
				}
			case string:
				var decordedObjString string
				if err := safecast.ToString(decordedObj, &decordedObjString); err != nil {
					t.Fatalf("failed to cast decoded object to string: %v", err)
				}
				if testObj != decordedObjString {
					t.Errorf("decoded data does not match original: expected %v, got %v", testObj, decordedObjString)
				}
			default:
				t.Fatalf("decoded object is not of type []byte: %T", testObj)
			}
		})
	}
}
