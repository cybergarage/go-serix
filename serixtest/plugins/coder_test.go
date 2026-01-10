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

package plugins

import (
	"testing"

	"github.com/cybergarage/go-serix/serix/document"
	"github.com/cybergarage/go-serix/serix/plugins"
	"github.com/cybergarage/go-serix/serixtest"
)

func TestPlugins(t *testing.T) {
	t.Helper()

	mgr := plugins.NewManager()

	t.Run("key", func(t *testing.T) {
		for _, coder := range mgr.KeyCoders() {
			t.Run(coder.Name(), func(t *testing.T) {
				serixtest.KeyCoderSuite(t, coder)
			})
		}
	})

	t.Run("object", func(t *testing.T) {
		objSerializer := []document.ObjectCoder{}
		objCompressor := []document.ObjectCoder{}
		for _, coder := range mgr.ObjectCoders() {
			switch coder.Type() {
			case document.ObjectSerializer:
				objSerializer = append(objSerializer, coder)
			case document.ObjectCompressor:
				objCompressor = append(objCompressor, coder)
			}
		}

		t.Run("serializer", func(t *testing.T) {
			for _, serializer := range objSerializer {
				t.Run(serializer.Name(), func(t *testing.T) {
					serixtest.ObjectSerializerSuite(t, serializer)
				})
			}
		})

		t.Run("compressor", func(t *testing.T) {
			for _, compressor := range objCompressor {
				t.Run(compressor.Name(), func(t *testing.T) {
					serixtest.ObjectCompressorSuite(t, compressor)
					for _, serializer := range objSerializer {
						t.Run(serializer.Name(), func(t *testing.T) {
							coder := document.NewChainCorder(serializer, compressor)
							serixtest.ObjectSerializerSuite(t, coder)
						})
					}
				})
			}
		})
	})
}
