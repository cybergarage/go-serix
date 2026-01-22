# Getting Started with go-serix

`go-serix` is a Go library that provides an idiomatic, schema-driven API for encoding and decoding structured data. It is designed to be lightweight and integration-friendly, with pluggable serialization formats and composable compression.

- Repository: https://github.com/cybergarage/go-serix
- Go reference: https://pkg.go.dev/github.com/cybergarage/go-serix

## Requirements

- Go (a recent version is recommended)
- A Go module (`go mod init ...`) for your project

## Installation

Add the module to your project:

```bash
go get github.com/cybergarage/go-serix@latest
```

## Key Concepts (as implemented in this repo)

### Document objects

In `go-serix`, a *document object* is represented by `document.Object` and encoded/decoded by an `document.ObjectCoder`.

An `ObjectCoder` is a simple interface with:

- `EncodeObject(w io.Writer, obj document.Object) error`
- `DecodeObject(r io.Reader) (document.Object, error)`
- `Name()` and `Type()` metadata

(See `serix/document/object_coder.go`.)

### Keys

A document key is represented by `document.Key` (a `[]any`) and can be created with:

- `document.NewKey()`
- `document.NewKeyWith(elems ...any)`

(See `serix/document/key.go`.)

### Plugins (coders)

The repository includes a plugin manager that registers built-in coders:

- Object serializers: `cbor`, `gob`, `json`
- Object compressors: `gzip`, `zlib`
- Key coder: `composite`

(See `serix/plugins/manager_impl.go`.)

## Quickstart: Encode / Decode with JSON

The JSON object coder is implemented in:

- `serix/plugins/document/object/json/coder.go`

It encodes/decodes `document.Object` using `encoding/json`.

### Example

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
	jsoncoder "github.com/cybergarage/go-serix/serix/plugins/document/object/json"
)

func main() {
	// document.Object is flexible; it can carry common Go values.
	// A typical choice is a map[string]any or []any.
	obj := document.Object(map[string]any{
		"hello": "world",
		"n":     42,
		"tags":  []any{"a", "b"},
	})

	coder := jsoncoder.NewCoder()

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

## Using the Plugin Manager

If you want to work with built-in coders dynamically, create a plugin manager:

```go
package main

import (
	"fmt"

	"github.com/cybergarage/go-serix/serix/plugins"
)

func main() {
	m := plugins.NewManager()

	for _, c := range m.ObjectCoders() {
		fmt.Printf("object coder: name=%s type=%v\n", c.Name(), c.Type())
	}

	for _, c := range m.KeyCoders() {
		fmt.Printf("key coder: %T\n", c)
	}
}
```

This is useful when you want to:
- expose selectable wire formats (e.g., JSON vs CBOR),
- configure compression (gzip/zlib) based on environment,
- keep your application-level interface stable while swapping codecs.

## Composing Serialization + Compression

`go-serix` supports chaining coders via `document.NewChainCorder(...)` (see `serix/document/chain_corder.go`).

A typical pattern is:

1. serialize structured data to bytes (e.g., CBOR/JSON/GOB),
2. compress the serialized bytes (e.g., gzip/zlib).

`NewChainCorder` runs `EncodeObject` left-to-right, and `DecodeObject` right-to-left.

### Example: JSON + gzip

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/go-serix/serix/document"
	gzipcoder "github.com/cybergarage/go-serix/serix/plugins/document/object/gzip"
	jsoncoder "github.com/cybergarage/go-serix/serix/plugins/document/object/json"
)

func main() {
	obj := document.Object(map[string]any{"hello": "world"})

	chain := document.NewChainCorder(
		jsoncoder.NewCoder(), // serializer (produces JSON bytes)
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

## Working with Keys

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

Notes:
- `Compare` uses `github.com/cybergarage/go-safecast/safecast` internally (see `serix/document/key.go`).
- Keep key element types consistent across your application for predictable ordering/comparison.

## Troubleshooting & Tips

### “unexpected type …” while decoding a chain
`NewChainCorder.DecodeObject` expects intermediate objects to be `[]byte` or `string` so it can feed them into the next decoder (see `serix/document/chain_corder.go`).

If you chain coders, make sure:
- serializers output `[]byte`/`string` (JSON/CBOR/GOB do),
- compressors accept `[]byte` (gzip/zlib require `[]byte` or `string`).

### gzip/zlib only compress byte-like objects
The gzip coder explicitly requires `[]byte` (and treats `string` by converting to bytes) (see `serix/plugins/document/object/gzip/coder.go`).
So you should always place compression *after* serialization.

## Next Steps

- Read the GoDoc: https://pkg.go.dev/github.com/cybergarage/go-serix
- Explore built-in coders under: `serix/plugins/document/object/`
- Decide on:
  - your default wire format (JSON / CBOR / GOB),
  - whether to add compression (gzip / zlib),
  - and whether you need dynamic selection via the plugin manager.