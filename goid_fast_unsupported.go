// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goidfast && !amd64 && !arm64

package runtime

// Goid returns the ID of the current goroutine.
//
// The goidfast fast path is only implemented for amd64 and arm64; on other
// architectures this build uses the portable slow path.
func Goid() uint64 { return goidSlow() }
