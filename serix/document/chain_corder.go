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
	"bytes"
	"fmt"
	"io"
	"strings"
)

type chainCorder struct {
	coders []ObjectCoder
}

// NewChainCorder creates a new chain coder with the specified coders.
func NewChainCorder(coders ...ObjectCoder) ObjectCoder {
	return &chainCorder{
		coders: coders,
	}
}

// Name returns the name of the coder.
func (mc *chainCorder) Name() string {
	names := make([]string, len(mc.coders))
	for i, coder := range mc.coders {
		names[i] = coder.Name()
	}
	return "multi(" + strings.Join(names, ",") + ")"
}

// Type returns the type of the coder.
func (mc *chainCorder) Type() CoderType {
	return ObjectSerializer | ObjectCompressor
}

func coderError(coder ObjectCoder, err error) error {
	return fmt.Errorf("%w (%s)", err, coder.Name())
}

// EncodeObject writes the specified object to the specified writer.
func (mc *chainCorder) EncodeObject(w io.Writer, obj Object) error {
	nextObject := obj
	for _, coder := range mc.coders {
		var buf bytes.Buffer
		if err := coder.EncodeObject(&buf, nextObject); err != nil {
			return coderError(coder, err)
		}
		// Make a copy to avoid issues with buffer reuse
		nextObject = append([]byte(nil), buf.Bytes()...)
	}
	// nextObject is now []byte, write it to the output writer
	if data, ok := nextObject.([]byte); ok {
		_, err := w.Write(data)
		return err
	}
	return nil
}

// DecodeObject returns the decorded object from the specified reader if available, otherwise returns an error.
func (mc *chainCorder) DecodeObject(r io.Reader) (Object, error) {
	nextReader := r
	var lastObject Object
	for i := len(mc.coders) - 1; i >= 0; i-- {
		coder := mc.coders[i]
		obj, err := coder.DecodeObject(nextReader)
		if err != nil {
			return nil, coderError(coder, err)
		}
		lastObject = obj
		if i == 0 {
			break
		}
		switch v := obj.(type) {
		case string:
			nextReader = strings.NewReader(v)
		case []byte:
			nextReader = bytes.NewReader(v)
		default:
			return nil, coderError(coder, fmt.Errorf("unexpected type %T", obj))
		}
	}
	return lastObject, nil
}
