// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime

import "runtime"

// goidSlow returns the id of the current goroutine by parsing the header
// of runtime.Stack. It is portable and always correct, and serves both as
// the default Goid implementation and as the fallback for the fast path.
//
// The format is stable: runtime.Stack writes a header like
//
//	goroutine 18446744073709551615 [running]:
//	golang.design/x/runtime.Goid...
//
// where the "goroutine " prefix is 10 bytes and the id is at most 20
// digits, so 32 bytes always holds "goroutine <id> ". The scan is bounded
// by the number of bytes Stack actually wrote so a missing trailing space
// (a future format change or a truncated write) cannot run past it.
func goidSlow() (id uint64) {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)
	for i := 10; i < n && buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}
	return id
}
