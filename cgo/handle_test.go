// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package cgo

import "testing"

func countHandles() int {
	siz := 0
	handles.Range(func(k, v interface{}) bool {
		siz++
		return true
	})
	return siz
}

func TestValueHandle(t *testing.T) {
	v := 42

	h1 := OpenHandle(v)
	h2 := OpenHandle(v)

	if uintptr(h1) == uintptr(h2) {
		t.Fatalf("duplicated Go values should have different handles")
	}

	h1v := h1.Value().(int)
	h2v := h2.Value().(int)
	if h1v != h2v {
		t.Fatalf("the Value of duplicated Go values are different: want %d, got %d", h1v, h2v)
	}
	if h1v != v {
		t.Fatalf("the Value of a handle does not match origin: want %v, got %v", v, h1v)
	}

	h1.Delete()
	h2.Delete()

	if siz := countHandles(); siz != 0 {
		t.Fatalf("handles are not deleted, want: %d, got %d", 0, siz)
	}
}

func TestPointerHandle(t *testing.T) {
	v := 42

	p1 := &v
	p2 := &v

	h1 := OpenHandle(p1)
	h2 := OpenHandle(p2)

	// Each OpenHandle call returns a distinct handle, even for pointers
	// that refer to the same object.
	if uintptr(h1) == uintptr(h2) {
		t.Fatalf("each OpenHandle call should return a distinct handle")
	}

	if h1.Value().(*int) != p1 {
		t.Fatalf("the Value of a handle does not match origin: want %p", p1)
	}
	if h2.Value().(*int) != p2 {
		t.Fatalf("the Value of a handle does not match origin: want %p", p2)
	}

	// Deleting one handle must not invalidate the other.
	h1.Delete()
	if _, ok := handles.Load(uintptr(h2)); !ok {
		t.Fatalf("deleting one handle wrongly invalidated another")
	}
	h2.Delete()

	if siz := countHandles(); siz != 0 {
		t.Fatalf("handles are not deleted: want %d, got %d", 0, siz)
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Fatalf("double Delete on a same handle did not trigger a panic")
	}()

	h1.Delete()
}

func TestNilHandle(t *testing.T) {
	var v *int

	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Fatalf("nil should not be created as a handle successfully")
	}()

	_ = OpenHandle(v)
}

// emptyA and emptyB are distinct zero-size types. Pointers to values of
// such types share the single runtime.zerobase address, so an identity
// scheme derived from the pointee address would collide them.
type emptyA struct{}
type emptyB struct{}

func TestZeroSizeHandle(t *testing.T) {
	a := &emptyA{}
	b := &emptyB{}

	ha := OpenHandle(a)
	hb := OpenHandle(b)

	if uintptr(ha) == uintptr(hb) {
		t.Fatalf("distinct zero-size pointers must get distinct handles")
	}

	if _, ok := ha.Value().(*emptyA); !ok {
		t.Fatalf("Value returned the wrong object for handle a: got %T", ha.Value())
	}
	if _, ok := hb.Value().(*emptyB); !ok {
		t.Fatalf("Value returned the wrong object for handle b: got %T", hb.Value())
	}

	// Deleting one must not invalidate (or panic on) the other.
	ha.Delete()
	if _, ok := handles.Load(uintptr(hb)); !ok {
		t.Fatalf("deleting handle a wrongly affected handle b")
	}
	hb.Delete()

	if siz := countHandles(); siz != 0 {
		t.Fatalf("handles are not deleted: want %d, got %d", 0, siz)
	}
}

func f1() {}
func f2() {}

type foo struct{}

func (f *foo) bar() {}
func (f *foo) wow() {}

func TestFuncHandle(t *testing.T) {
	h1 := OpenHandle(f1)
	h2 := OpenHandle(f2)
	h3 := OpenHandle(f2)

	// Every OpenHandle call returns a distinct handle, including for the
	// same function value.
	if h1 == h2 || h2 == h3 || h1 == h3 {
		t.Fatalf("each OpenHandle call should return a distinct handle")
	}

	f := foo{}
	h4 := OpenHandle(f.bar)
	h5 := OpenHandle(f.bar)
	h6 := OpenHandle(f.wow)

	if h4 == h5 || h5 == h6 || h4 == h6 {
		t.Fatalf("each OpenHandle call should return a distinct handle")
	}

	for _, h := range []Handle{h1, h2, h3, h4, h5, h6} {
		h.Delete()
	}
	if siz := countHandles(); siz != 0 {
		t.Fatalf("handles are not deleted: want %d, got %d", 0, siz)
	}
}

func BenchmarkHandle(b *testing.B) {
	b.Run("non-concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := OpenHandle(i)
			_ = h.Value()
			h.Delete()
		}
	})
	b.Run("concurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var v int
			for pb.Next() {
				h := OpenHandle(v)
				_ = h.Value()
				h.Delete()
			}
		})
	})
}
