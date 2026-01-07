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

package json

import (
	"encoding/json"
	"io"

	"github.com/cybergarage/go-serix/serix/document"
)

// Coder represents a JSON serializer.
type Coder struct {
}

// NewCoder returns a new JSON serializer instance.
func NewCoder() *Coder {
	return &Coder{}
}

// Name returns the name of the coder.
func (s *Coder) Name() string {
	return "json"
}

// EncodeObject writes the specified object to the specified writer.
func (s *Coder) EncodeObject(w io.Writer, obj document.Object) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

// DecodeObject returns the decoded object from the specified reader if available, otherwise returns an error.
func (s *Coder) DecodeObject(r io.Reader) (document.Object, error) {
	var obj document.Object
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
