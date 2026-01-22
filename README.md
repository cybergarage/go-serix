# go-serix
![](https://img.shields.io/badge/status-Work%20In%20Progress-8A2BE2)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-serix)
[![test](https://github.com/cybergarage/go-serix/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-serix/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-serix.svg)](https://pkg.go.dev/github.com/cybergarage/go-serix)
[![codecov](https://codecov.io/gh/cybergarage/go-serix/graph/badge.svg?token=GOLCBMUVB1)](https://codecov.io/gh/cybergarage/go-serix)

## Overview

`go-serix` is a Go library that provides an idiomatic, schema-driven API for encoding and decoding structured data. Serix is a coined term used in this project to describe a compact, schema-first serialization approach for representing and exchanging structured data. The library is designed to be lightweight and easy to integrate into existing Go services, with a focus on clarity, maintainability, and performance.

## Features

- Stable high-level interfaces: keep your application-facing API the same
- Pluggable serialization formats: switch the underlying wire format without changing the upper-layer code
- Composable compression: combine your chosen serialization format with different compression strategies
- Integration-friendly design for production Go services

The name "Serix" is derived from "Serialize" and "X" (standing for extensibility), reflecting the project's goal of providing a flexible and extensible serialization framework.

## Getting Started with go-serix

`go-serix` is a Go library that provides an idiomatic, schema-driven API for encoding and decoding structured data. It is designed to be lightweight and integration-friendly, with pluggable serialization formats and composable compression.

- Repository: https://github.com/cybergarage/go-serix
- Go reference: https://pkg.go.dev/github.com/cybergarage/go-serix

### Installation

Add the module to your project:

```bash
go get github.com/cybergarage/go-serix@latest
```

### Key Concepts (as implemented in this repo)

#### Document objects

In `go-serix`, a *document object* is represented by `document.Object` and encoded/decoded by an `document.ObjectCoder`.

An `ObjectCoder` is a simple interface with:

- `EncodeObject(w io.Writer, obj document.Object) error`
- `DecodeObject(r io.Reader) (document.Object, error)`
- `Name()` and `Type()` metadata

(See `serix/document/object_coder.go`.)

#### Keys

A document key is represented by `document.Key` (a `[]any`) and can be created with:

- `document.NewKey()`
- `document.NewKeyWith(elems ...any)`

(See `serix/document/key.go`.)

#### Plugins (coders)

The repository includes a plugin manager that registers built-in coders:

- Object serializers: `cbor`, `gob`, `cbor`
- Object compressors: `gzip`, `zlib`
- Key coder: `composite`

### Quickstart: Encode / Decode with CBOR

The CBOR object coder is implemented in:

- `serix/plugins/document/object/cbor/coder.go`

It encodes/decodes `document.Object` using `encoding/cbor`.

#### Example

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
	cborcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/cbor"
)

func main() {
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
```

### Composing Serialization + Compression

`go-serix` supports chaining coders via `document.NewChainCorder(...)` (see `serix/document/chain_corder.go`).

A typical pattern is:

1. serialize structured data to bytes (e.g., CBOR/CBOR/GOB),
2. compress the serialized bytes (e.g., gzip/zlib).

`NewChainCorder` runs `EncodeObject` left-to-right, and `DecodeObject` right-to-left.

#### Example: CBOR + gzip

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
	gzipcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/gzip"
	cborcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/cbor"
)

func main() {
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
```

### Working with Keys

Keys are `[]any` values used to identify document objects:

```go
package main

import (
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
)

func main() {
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
```