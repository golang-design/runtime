// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build linux
// +build linux

package runtime

import (
	"fmt"
	"os"
)

// numThreads returns the number of OS threads in the current process by
// counting the per-thread directories under /proc/self/task.
//
// This replaces the previous `ps hH | wc -l` implementation, which forked
// a shell on every call (milliseconds, plus a hard dependency on bash and
// ps that silently fails in minimal containers). Reading /proc is orders
// of magnitude faster and works without any external command.
func numThreads() int {
	f, err := os.Open("/proc/self/task")
	if err != nil {
		if debug {
			fmt.Printf("mkill: failed to open /proc/self/task: %v\n", err)
		}
		return 0
	}
	defer f.Close()

	// Readdirnames avoids an Lstat per entry; each entry is one thread.
	names, err := f.Readdirnames(-1)
	if err != nil {
		if debug {
			fmt.Printf("mkill: failed to read /proc/self/task: %v\n", err)
		}
		return 0
	}
	return len(names)
}
