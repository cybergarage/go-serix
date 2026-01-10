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

// KeyDecoder decodes the specified bytes into a key.
type KeyDecoder interface {
	// DecodeKey returns the decoded key from the specified bytes if available, otherwise returns an error.
	DecodeKey([]byte) (Key, error)
}

// KeyEncoder encodes the specified key into bytes.
type KeyEncoder interface {
	// EncodeKey returns the encoded bytes from the specified key if available, otherwise returns an error.
	EncodeKey(Key) ([]byte, error)
}

// A KeyCoder includes key decoder and encoder interfaces.
type KeyCoder interface {
	// Name returns the name of the coder.
	Name() string
	// KeyDecoder decodes the specified bytes into a key.
	KeyDecoder
	// KeyEncoder encodes the specified key into bytes.
	KeyEncoder
}
