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
	_ "embed"
	"testing"

	"github.com/cybergarage/go-serix/serix/document"
)

// ObjectCoderSuite tests the specified document coder.
func ObjectCoderSuite(t *testing.T, coder document.Coder) {
	t.Helper()

	testFuncs := []struct {
		name string
		fn   func(*testing.T, document.Coder)
	}{
		{"primitive", primitiveDocumentTest},
		{"array", arrayDocumentTest},
		{"map", mapDocumentTest},
	}

	for _, testFunc := range testFuncs {
		t.Run(testFunc.name, func(t *testing.T) {
			testFunc.fn(t, coder)
		})
	}
}
