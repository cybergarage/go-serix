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

package document

import (
	"io"
)

// ObjectDecoder represets an interface for decoding objects from an input stream.
type ObjectDecoder interface {
	// DecodeObject reads an object from the specified reader.
	DecodeObject(r io.Reader) (Object, error)
}

// ObjectEncoder represets an interface for encoding objects to an output stream.
type ObjectEncoder interface {
	// EncodeObject writes the specified object to the specified writer.
	EncodeObject(w io.Writer, obj Object) error
}

// Coder represents an interface for encoding and decoding objects.
type Coder interface {
	// Name returns the name of the coder.
	Name() string
	// Type returns the type of the coder.
	Type() CoderType
	// ObjectDecoder returns the object decoder.
	ObjectDecoder
	// ObjectEncoder returns the object encoder.
	ObjectEncoder
}
