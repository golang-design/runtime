// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build !goidfast

package runtime

// Goid returns the ID of the current goroutine.
//
// By default it parses runtime.Stack, which is portable and safe but costs
// on the order of microseconds. Build with -tags goidfast on amd64 or
// arm64 for a ~nanosecond implementation that reads the id directly from
// the runtime's goroutine structure (see goid_fast.go); that build falls
// back to this implementation if it cannot calibrate.
func Goid() uint64 { return goidSlow() }
