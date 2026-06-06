// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows
// +build windows

package runtime

import (
	"fmt"
	"syscall"
	"unsafe"
)

// The standard syscall package wraps CreateToolhelp32Snapshot and
// CloseHandle but not Thread32First/Thread32Next, so bind those (and the
// THREADENTRY32 layout) directly to stay dependency-free.
var (
	kernel32      = syscall.NewLazyDLL("kernel32.dll")
	thread32First = kernel32.NewProc("Thread32First")
	thread32Next  = kernel32.NewProc("Thread32Next")
)

// threadEntry32 mirrors the Win32 THREADENTRY32 structure.
type threadEntry32 struct {
	Size           uint32
	Usage          uint32
	ThreadID       uint32
	OwnerProcessID uint32
	BasePri        int32
	DeltaPri       int32
	Flags          uint32
}

// checkwork verifies the toolhelp thread snapshot is usable before the
// watcher starts. Unlike the Unix builds, Windows counts threads through
// a syscall rather than an external command, so there is no bash/ps
// dependency to validate here.
func checkwork() error {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return fmt.Errorf("runtime: failed to use the package: %w", err)
	}
	syscall.CloseHandle(snapshot)
	return nil
}

// numThreads returns the number of OS threads owned by the current
// process. The previous implementation returned GetCurrentProcessId,
// i.e. the PID, which is not a thread count and made the whole
// max-thread mechanism meaningless on Windows.
func numThreads() int {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		if debug {
			fmt.Printf("mkill: failed to snapshot threads: %v\n", err)
		}
		return 0
	}
	defer syscall.CloseHandle(snapshot)

	owner := uint32(pid)
	var entry threadEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	ret, _, err := thread32First.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		if debug {
			fmt.Printf("mkill: failed to walk threads: %v\n", err)
		}
		return 0
	}

	count := 0
	for {
		if entry.OwnerProcessID == owner {
			count++
		}
		ret, _, _ := thread32Next.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			// ERROR_NO_MORE_FILES marks the end of the snapshot.
			break
		}
	}
	return count
}
