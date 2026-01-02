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
	"testing"

	"github.com/cybergarage/go-serix/serix/document"
)

// KeyCoderTest tests the encoding and decoding of keys using the provided KeyCoder.
func KeyCoderTest(t *testing.T, coder document.KeyCoder) {
	t.Helper()

	tests := []struct {
		name string
		test func(t *testing.T, coder document.KeyCoder)
	}{
		{
			name: "RoundTripKeyTest",
			test: RoundTripKeyTest,
		},
		{
			name: "SortableKeyTest",
			test: SortableKeyTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, coder)
		})
	}
}
