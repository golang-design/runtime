// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goidfast && (amd64 || arm64)

package runtime

import (
	"sync"
	"unsafe"
)

// getg returns a pointer to the current goroutine's runtime.g structure.
// It is implemented in assembly per architecture (getg_<arch>.s).
func getg() unsafe.Pointer

var (
	goidOffset     uintptr
	goidCalibrated bool
)

func init() { calibrateGoid() }

// calibrateGoid discovers the byte offset of the int64 goid field inside
// runtime.g, instead of hard-coding it (which would break on every Go
// release that reorders the structure). It samples two goroutines, each
// reporting its slow-path id and its g pointer, then finds the 8-aligned
// offset whose word equals the known id in BOTH samples. Requiring two
// distinct ids to agree rules out a coincidental match against some other
// field. If no offset matches, the fast path stays disabled and Goid falls
// back to the slow path.
func calibrateGoid() {
	type sample struct {
		id uint64
		g  unsafe.Pointer
	}
	ch := make(chan sample, 2)
	var start sync.WaitGroup
	start.Add(1)
	for i := 0; i < 2; i++ {
		go func() {
			start.Wait() // ensure two distinct, live goroutines
			ch <- sample{id: goidSlow(), g: getg()}
		}()
	}
	start.Done()
	s1, s2 := <-ch, <-ch

	// runtime.g is well over 512 bytes; the goid field has always lived
	// near the front. Scanning 8-aligned words keeps reads in-bounds.
	const maxScan = 512
	for off := uintptr(0); off+8 <= maxScan; off += 8 {
		if wordAt(s1.g, off) == s1.id && wordAt(s2.g, off) == s2.id {
			goidOffset = off
			goidCalibrated = true
			return
		}
	}
}

func wordAt(p unsafe.Pointer, off uintptr) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(p) + off))
}

// Goid returns the ID of the current goroutine.
//
// This is the -tags goidfast build: it reads the id directly from the
// runtime goroutine structure (~nanoseconds). If offset calibration failed
// at init, it falls back to the portable slow path.
func Goid() uint64 {
	if goidCalibrated {
		return wordAt(getg(), goidOffset)
	}
	return goidSlow()
}
