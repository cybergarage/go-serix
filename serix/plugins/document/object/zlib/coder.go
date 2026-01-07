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

package zlib

import (
	"compress/zlib"
	"fmt"
	"io"

	"github.com/cybergarage/go-serix/serix/document"
)

// Coder implements document compression and decompression using zlib.
type Coder struct{}

// NewCoder creates a new zlib coder instance.
func NewCoder() *Coder {
	return &Coder{}
}

// Name returns the name of the coder.
func (s *Coder) Name() string {
	return "zlib"
}

// Type returns the type of the coder.
func (s *Coder) Type() document.CoderType {
	return document.ObjectCompressor
}

// EncodeObject writes the specified object to the specified writer.
func (s *Coder) EncodeObject(w io.Writer, obj document.Object) error {
	var data []byte
	switch v := obj.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("zlib coder only supports []byte and string objects")
	}

	zw := zlib.NewWriter(w)
	defer zw.Close()

	_, err := zw.Write(data)
	return err
}

// DecodeObject returns the decoded object from the specified reader if available, otherwise returns an error.
func (s *Coder) DecodeObject(r io.Reader) (document.Object, error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	buf, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
