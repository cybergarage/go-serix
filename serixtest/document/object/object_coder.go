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

package object

import (
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/cybergarage/go-pict/pict"
	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-serix/serix/document"
)

//go:embed go_types.pict
var goTypes []byte

func deepEqual(x, y any) error {
	if reflect.DeepEqual(x, y) {
		return nil
	}
	if safecast.Equal(x, y) {
		return nil
	}
	return fmt.Errorf("%v (%T) != %v (%T)", x, x, y, y)
}

func primitiveDocumentTest(t *testing.T, coder document.Coder) {
	t.Helper()

	pict := pict.NewParserWithBytes(goTypes)
	err := pict.Parse()
	if err != nil {
		t.Fatal(err)
	}

	pictParams := pict.Params()
	for n, pictParam := range pictParams {
		t.Run(string(pictParam), func(t *testing.T) {
			for _, pictCase := range pict.Cases() {
				pictElem := pictCase[n]
				pictType, err := pictParam.Type()
				if err != nil {
					t.Error(err)
					return
				}
				obj, err := pictElem.CastTo(pictType)
				if err != nil {
					t.Error(err)
					return
				}

				var w bytes.Buffer
				err = coder.EncodeDocument(&w, obj)
				if err != nil {
					t.Error(err)
					return
				}

				r := bytes.NewReader(w.Bytes())
				decObj, err := coder.DecodeDocument(r)
				if err != nil {
					t.Error(err)
					return
				}

				err = deepEqual(decObj, obj)
				if err != nil {
					t.Error(err)
					return
				}
			}
		})
	}
}

func mapDocumentTest(t *testing.T, coder document.Coder) {
	t.Helper()

	pict := pict.NewParserWithBytes(goTypes)
	err := pict.Parse()
	if err != nil {
		t.Fatal(err)
	}

	pictParams := pict.Params()
	for _, pictCase := range pict.Cases() {
		obj := map[string]any{}
		for n, pictParam := range pictParams {
			name := string(pictParam)
			pictElem := pictCase[n]
			pictType, err := pictParam.Type()
			if err != nil {
				t.Error(err)
				return
			}
			v, err := pictElem.CastTo(pictType)
			if err != nil {
				t.Error(err)
				return
			}
			obj[name] = v
		}

		var w bytes.Buffer
		err = coder.EncodeDocument(&w, obj)
		if err != nil {
			t.Error(err)
			return
		}

		r := bytes.NewReader(w.Bytes())
		decObj, err := coder.DecodeDocument(r)
		if err != nil {
			t.Error(err)
			return
		}

		err = deepEqual(decObj, obj)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

// nolint:gosec
func arrayDocumentTest(t *testing.T, coder document.Coder) {
	t.Helper()

	pict := pict.NewParserWithBytes(goTypes)
	err := pict.Parse()
	if err != nil {
		t.Fatal(err)
	}

	shuffleArray := func(array []any) {
		n := len(array)
		for i := n - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			array[i], array[j] = array[j], array[i]
		}
	}

	pictParams := pict.Params()
	for _, pictCase := range pict.Cases() {
		obj := []any{}
		for n, pictParam := range pictParams {
			pictElem := pictCase[n]
			pictType, err := pictParam.Type()
			if err != nil {
				t.Error(err)
				return
			}
			v, err := pictElem.CastTo(pictType)
			if err != nil {
				t.Error(err)
				return
			}
			obj = append(obj, v)
		}

		// Non-shuffled array

		var w bytes.Buffer
		err = coder.EncodeDocument(&w, obj)
		if err != nil {
			t.Error(err)
			return
		}

		r := bytes.NewReader(w.Bytes())
		decObj, err := coder.DecodeDocument(r)
		if err != nil {
			t.Error(err)
			return
		}

		err = deepEqual(decObj, obj)
		if err != nil {
			t.Error(err)
			return
		}

		// Shuffled array

		w.Reset()

		shuffleArray(obj)

		err = coder.EncodeDocument(&w, obj)
		if err != nil {
			t.Error(err)
			return
		}

		r = bytes.NewReader(w.Bytes())
		decObj, err = coder.DecodeDocument(r)
		if err != nil {
			t.Error(err)
			return
		}

		err = deepEqual(decObj, obj)
		if err != nil {
			t.Error(err)
			return
		}

		// Random reduced array

		w.Reset()

		nObj := rand.Intn(len(obj)-1) + 1
		obj = obj[:nObj]

		err = coder.EncodeDocument(&w, obj)
		if err != nil {
			t.Error(err)
			return
		}

		b := w.Bytes()
		r = bytes.NewReader(b)
		decObj, err = coder.DecodeDocument(r)
		if err != nil {
			t.Error(err)
			return
		}

		err = deepEqual(decObj, obj)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
