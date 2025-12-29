package document

import (
	"testing"
)

func TestKeyCompare(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		cases := []struct {
			a Key
			b Key
		}{
			{NewKeyWith(), NewKeyWith()},
			{NewKeyWith(1, "a", true), NewKeyWith(1, "a", true)},
			{NewKeyWith(int(1), int64(1)), NewKeyWith(int8(1), uint(1))},
			{NewKeyWith("foo"), NewKeyWith("foo")},
			{NewKeyWith(uint32(10), float64(10)), NewKeyWith(int16(10), float32(10))},
		}

		for _, c := range cases {
			if cmp, err := c.a.Compare(c.b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if cmp != 0 {
				t.Fatalf("expected 0, got %d for %v vs %v", cmp, c.a, c.b)
			}
			if cmp, err := c.b.Compare(c.a); err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if cmp != 0 {
				t.Fatalf("expected 0, got %d for %v vs %v", cmp, c.b, c.a)
			}
		}
	})

	t.Run("order", func(t *testing.T) {
		// Differ at first element
		a := NewKeyWith(1, "a")
		b := NewKeyWith(2, "a")
		if cmp, err := a.Compare(b); err != nil || cmp >= 0 {
			t.Fatalf("expected a<b, cmp=%d err=%v", cmp, err)
		}
		if cmp, err := b.Compare(a); err != nil || cmp <= 0 {
			t.Fatalf("expected b>a, cmp=%d err=%v", cmp, err)
		}

		// Differ at later element
		c := NewKeyWith("x", 10)
		d := NewKeyWith("x", 20)
		if cmp, err := c.Compare(d); err != nil || cmp >= 0 {
			t.Fatalf("expected c<d, cmp=%d err=%v", cmp, err)
		}
		if cmp, err := d.Compare(c); err != nil || cmp <= 0 {
			t.Fatalf("expected d>c, cmp=%d err=%v", cmp, err)
		}

		// Prefix ordering: shorter prefix should be less when all common elements equal
		e := NewKeyWith("p")
		f := NewKeyWith("p", 1)
		if cmp, err := e.Compare(f); err != nil || cmp >= 0 {
			t.Fatalf("expected e<f (prefix shorter), cmp=%d err=%v", cmp, err)
		}
		if cmp, err := f.Compare(e); err != nil || cmp <= 0 {
			t.Fatalf("expected f>e (prefix longer), cmp=%d err=%v", cmp, err)
		}
	})

	t.Run("mixed-numeric", func(t *testing.T) {
		a := NewKeyWith(int8(5))
		b := NewKeyWith(uint16(7))
		if cmp, err := a.Compare(b); err != nil || cmp >= 0 {
			t.Fatalf("expected 5<7, cmp=%d err=%v", cmp, err)
		}

		c := NewKeyWith(float32(10.0))
		d := NewKeyWith(int64(10))
		if cmp, err := c.Compare(d); err != nil || cmp != 0 {
			t.Fatalf("expected 10.0==10, cmp=%d err=%v", cmp, err)
		}
	})

	t.Run("incomparable-types", func(t *testing.T) {
		// Some pairs should be incomparable for safecast (e.g., map vs number).
		a := NewKeyWith(map[string]int{"a": 1})
		b := NewKeyWith(1)
		if _, err := a.Compare(b); err == nil {
			t.Fatalf("expected error comparing map and number")
		}
	})
}
