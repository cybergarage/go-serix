// Copyright (C) 2025 The go-serix Authors.
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

package serix_test

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"

	cborcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/cbor"
	gzipcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/gzip"
)

func Example_objectCoder() {
	// document.Object is flexible; it can carry common Go values.
	// A typical choice is a map[string]any or []any.
	obj := document.Object(map[string]any{
		"hello": "world",
		"n":     42,
		"tags":  []any{"a", "b"},
	})

	coder := cborcoder.NewCoder()

	// Encode
	var buf bytes.Buffer
	if err := coder.EncodeObject(&buf, obj); err != nil {
		panic(err)
	}

	// Decode
	decoded, err := coder.DecodeObject(&buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("encoded bytes: %s\n", buf.String())
	fmt.Printf("decoded value: %#v\n", decoded)
}

func Example_chainCorder() {
	obj := document.Object(map[string]any{"hello": "world"})

	chain := document.NewChainCorder(
		cborcoder.NewCoder(), // serializer (produces CBOR bytes)
		gzipcoder.NewCoder(), // compressor (compresses []byte)
	)

	// Encode
	var buf bytes.Buffer
	if err := chain.EncodeObject(&buf, obj); err != nil {
		panic(err)
	}

	// Decode (returns the original structured object)
	decoded, err := chain.DecodeObject(bytes.NewReader(buf.Bytes()))
	if err != nil {
		panic(err)
	}

	fmt.Printf("compressed size: %d bytes\n", len(buf.Bytes()))
	fmt.Printf("decoded value: %#v\n", decoded)
}

func Example_keyCorder() {
	k1 := document.NewKeyWith("user", uint64(123))
	k2 := document.NewKeyWith("user", uint64(456))

	cmp, err := k1.Compare(k2)
	if err != nil {
		panic(err)
	}

	fmt.Println("k1:", k1.String())
	fmt.Println("k2:", k2.String())
	fmt.Println("compare:", cmp)
}
