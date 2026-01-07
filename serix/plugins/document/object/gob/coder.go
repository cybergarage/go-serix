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

package gob

import (
	"encoding/gob"
	"io"

	"github.com/cybergarage/go-serix/serix/document"
)

func init() {
	// Register commonly used composite types when carried in interface fields.
	gob.Register([]any{})
	gob.Register(map[string]any{})
}

// Coder represents a gob serializer.
type Coder struct{}

// anyWrapper is used to reliably encode/decode arbitrary values, including nil,
// by carrying them in a known interface-typed field for gob.
type anyWrapper struct {
	V any
}

// NewCoder returns a new gob serializer instance.
func NewCoder() *Coder {
	return &Coder{}
}

// Name returns the name of the coder.
func (s *Coder) Name() string {
	return "gob"
}

// EncodeObject writes the specified object to the specified writer.
func (s *Coder) EncodeObject(w io.Writer, obj document.Object) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(&anyWrapper{V: obj})
}

// DecodeObject returns the decoded object from the specified reader if available, otherwise returns an error.
func (s *Coder) DecodeObject(r io.Reader) (document.Object, error) {
	var wrap anyWrapper
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&wrap); err != nil {
		return nil, err
	}
	return wrap.V, nil
}
