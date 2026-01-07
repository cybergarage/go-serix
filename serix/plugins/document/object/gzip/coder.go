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

package gzip

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/cybergarage/go-serix/serix/document"
)

// Coder implements document compression and decompression using gzip.
type Coder struct{}

// NewCoder creates a new gzip coder instance.
func NewCoder() *Coder {
	return &Coder{}
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
		return fmt.Errorf("gzip coder only supports []byte object")
	}

	// Create gzip writer
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()

	// Write compressed data
	_, err := gzipWriter.Write(data)
	return err
}

// DecodeObject returns the decoded object from the specified reader if available, otherwise returns an error.
func (s *Coder) DecodeObject(r io.Reader) (document.Object, error) {
	// Create gzip reader
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Read decompressed data
	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}
