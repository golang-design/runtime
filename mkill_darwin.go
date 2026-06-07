// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin
// +build darwin

package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// numThreads counts the process's OS threads via `ps M | wc -l`.
//
// Unlike Linux (which reads /proc) and Windows (which uses a toolhelp
// snapshot), darwin has no cheap dependency-free thread count: /proc does
// not exist, and the libproc/Mach alternatives require cgo or
// Apple-discouraged raw syscalls. macOS always ships ps and a shell and is
// not a minimal-container target, so the shell-out is kept here.
var cmdThreads = fmt.Sprintf("ps M %d | wc -l", pid)

func numThreads() int {
	out, err := exec.Command("bash", "-c", cmdThreads).Output()
	if err != nil && debug {
		fmt.Printf("mkill: failed to fetch #threads: %v\n", err)
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil && debug {
		fmt.Printf("mkill: failed to parse #threads: %v\n", err)
		return 0
	}
	return n
}
