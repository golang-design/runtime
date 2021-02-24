// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows
// +build windows

package runtime

import "fmt"

var (
	kernel32 = syscall.NewLazyDLL("kernel32")
	getPid   = kernel32.NewProc("GetCurrentProcessId")
)

func numThreads() int {
	pid, _, err := getPid.Call()
	if pid == 0 {
		panic(fmt.Sprintf("failed to get process id: %v", err))
	}
	return int(pid)
}
