// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package cgo is an implementation of golang.org/issue/37033.
//
// See golang.org/cl/294670 for code review discussion.
package cgo // import "golang.design/x/runtime/cgo"

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Handle provides a safe representation to pass Go values between C and
// Go back and forth. The zero value of a handle is not a valid handle,
// and thus safe to use as a sentinel in C APIs.
//
// The underlying type of Handle may change, but the value is guaranteed
// to fit in an integer type that is large enough to hold the bit pattern
// of any pointer. For instance, on the Go side:
//
// 	package main
//
// 	/*
// 	extern void MyGoPrint(unsigned long long handle);
// 	void myprint(unsigned long long handle);
// 	*/
// 	import "C"
// 	import "golang.design/x/runtime/cgo"
//
// 	//export MyGoPrint
// 	func MyGoPrint(handle C.ulonglong) {
// 		h := cgo.Handle(handle)
// 		val := h.Value().(int)
// 		println(val)
// 		h.Delete()
// 	}
//
// 	func main() {
// 		val := 42
// 		C.myprint(C.ulonglong(cgo.OpenHandle(val)))
// 		// Output: 42
// 	}
//
// and on the C side:
//
// 	// A Go function
// 	extern void MyGoPrint(unsigned long long handle);
//
// 	// A C function
// 	void myprint(unsigned long long handle) {
// 	    MyGoPrint(handle);
// 	}
type Handle uintptr

// OpenHandle returns a handle for a given value. Each call to OpenHandle
// returns a distinct handle, even for the same value; in particular,
// pointers, slices, maps, channels, or functions that refer to the same
// object still receive distinct handles. A nil value (a nil pointer,
// slice, map, channel, or function) must not be used.
//
// The handle is valid until the program calls Delete on it. The handle
// uses resources, and this package assumes that C code may hold on to
// the handle, so a program must explicitly call Delete when the handle
// is no longer needed.
//
// The intended use is to pass the returned handle to C code, which
// passes it back to Go, which calls Value. See an example in the
// comments of the Handle definition.
func OpenHandle(v interface{}) Handle {
	// Reject nil reference types. Unlike the standard runtime/cgo,
	// this package documents and guarantees that a nil value is never
	// turned into a handle, so it can serve as a safe sentinel.
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Slice,
		reflect.Map, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			panic("cgo: cannot use Handle for nil value")
		}
	}

	h := atomic.AddUintptr(&handleIdx, 1)
	if h == 0 {
		panic("cgo: ran out of handle space")
	}

	handles.Store(h, v)
	return Handle(h)
}

// Delete invalidates a handle. This method must be called when C code no
// longer has a copy of the handle, and the program no longer needs the
// Go value that associated with the handle.
//
// The method panics if the handle is invalid already.
func (h Handle) Delete() {
	_, ok := handles.LoadAndDelete(uintptr(h))
	if !ok {
		panic("cgo: misuse of an invalid Handle")
	}
}

// Value returns the associated Go value for a valid handle.
//
// The method panics if the handle is invalid already.
func (h Handle) Value() interface{} {
	v, ok := handles.Load(uintptr(h))
	if !ok {
		panic("cgo: misuse of an invalid Handle")
	}
	return v
}

var (
	handles   = &sync.Map{} // map[uintptr]interface{}
	handleIdx uintptr       // accessed atomically
)
