// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime

// Goid returns the ID of current goroutine.
//
// This implementation based on the facts that runtime.Stack gives
// information like:
//
//   goroutine 18446744073709551615 [running]:
//   golang.design/x/runtime.Goid...
//
// This format stands for more than 10 years.
// Since commit 4dfd7fdde5957e4f3ba1a0285333f7c807c28f03, a goroutine id
// ends with a white space. Go 1 compatability promise garantees all
// versions of Go can use this function.
func Goid() (id uint64) {
	// The prefix "goroutine " is 10 bytes; a uint64 id is at most 20
	// digits, so 32 bytes always holds "goroutine <id> ". Bound the scan
	// by the number of bytes Stack actually wrote so a missing trailing
	// space (e.g. a future format change or a truncated write) cannot run
	// the loop past the written region.
	var buf [32]byte
	n := Stack(buf[:], false)
	for i := 10; i < n && buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}
	return id
}
