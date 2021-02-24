// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build linux
// +build linux

package runtime

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var cmdThreads = fmt.Sprintf("ps hH p %d | wc -l", pid)

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
