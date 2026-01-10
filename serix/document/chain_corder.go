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
	var lastWriter *bytes.Buffer
	for _, coder := range mc.coders {
		if lastWriter == nil {
			lastWriter = bytes.NewBuffer(nil)
		} else {
			lastWriter.Reset()
		}
		if err := coder.EncodeObject(lastWriter, nextObject); err != nil {
			return coderError(coder, err)
		}
		nextObject = lastWriter.Bytes()
	}
	if lastWriter == nil {
		return nil
	}
	_, err := w.Write(lastWriter.Bytes())
	return err
}

// DecodeObject returns the decorded object from the specified reader if available, otherwise returns an error.
func (mc *chainCorder) DecodeObject(r io.Reader) (Object, error) {
	lastReader := r
	var lastObject any
	var err error
	for i := len(mc.coders) - 1; i >= 0; i-- {
		coder := mc.coders[i]
		lastObject, err = coder.DecodeObject(lastReader)
		if err != nil {
			return nil, coderError(coder, err)
		}
		if i == 0 {
			break
		}
		nextObj, ok := lastObject.([]byte)
		if !ok {
			err := fmt.Errorf("expected []byte, got %T", lastObject)
			return nil, coderError(coder, err)
		}
		lastReader = bytes.NewReader(nextObj)
	}
	return lastObject, nil
}
