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
	"github.com/cybergarage/go-serix/serix/document"
	"github.com/cybergarage/go-serix/serix/plugins/document/key/composite"
	"github.com/cybergarage/go-serix/serix/plugins/document/object/cbor"
	"github.com/cybergarage/go-serix/serix/plugins/document/object/gob"
	"github.com/cybergarage/go-serix/serix/plugins/document/object/gzip"
	"github.com/cybergarage/go-serix/serix/plugins/document/object/json"
	"github.com/cybergarage/go-serix/serix/plugins/document/object/zlib"
)

type manager struct {
	keyCorders   []document.KeyCoder
	objectCoders []document.ObjectCoder
}

// NewManager returns a new plugin manager instance.
func NewManager() Manager {
	keyCoders := []document.KeyCoder{
		composite.NewCoder(),
	}
	objCoders := []document.ObjectCoder{
		cbor.NewCoder(),
		gob.NewCoder(),
		gzip.NewCoder(),
		json.NewCoder(),
		zlib.NewCoder(),
	}
	manager := &manager{
		keyCorders:   keyCoders,
		objectCoders: objCoders,
	}
	return manager
}

func (m *manager) KeyCoders() []document.KeyCoder {
	return m.keyCorders
}

func (m *manager) ObjectCoders() []document.ObjectCoder {
	return m.objectCoders
}
