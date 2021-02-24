// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	pid = os.Getpid()
	// minimum number of threads required by the runtime
	minThreads = int32(runtime.NumCPU()) + 2
	// 2 meaning runtime sysmon thread + template thread
	maxThreads = int32(runtime.NumCPU()) + 2
	interval   = time.Second
	debug      = false
)

// NumThreads returns the number of threads that currently exist.
func NumThreads() int {
	return numThreads()
}

var started int32

// SetMaxThreads sets the maximum number of system threads that allowed
// in a Go program and returns the previous setting. If n is lower than
// minimum required number of threads, it does not change the current
// setting. The minimum allowed number of threads of a program is
// runtime.NumCPU() + 2.
func SetMaxThreads(n int) int {
	finished := atomic.CompareAndSwapInt32(&started, 0, 1)
	if !finished {
		err := checkwork()
		if err != nil {
			return 0
		}
		if debug {
			fmt.Printf("runtime: pid %v, maxThread %v, interval %v\n",
				pid, maxThreads, interval)
		}

		wg := sync.WaitGroup{}
		go func() {
			t := time.NewTicker(interval)
			for {
				select {
				case <-t.C:
					n := NumThreads()
					nkill := int32(n) - atomic.LoadInt32(&maxThreads)
					if nkill <= 0 {
						if debug {
							fmt.Printf("runtime: checked #threads total %v / max %v\n",
								n, maxThreads)
						}
						continue
					}
					wg.Add(int(nkill))
					for i := int32(0); i < nkill; i++ {
						go func() {
							runtime.LockOSThread()
							wg.Done()
						}()
					}
					wg.Wait()
					if debug {
						fmt.Printf("runtime: killing #threads, remaining: %v\n", n)
					}
				}
			}
		}()
	}

	if n < int(minThreads) {
		return int(atomic.LoadInt32(&maxThreads))
	}

	return int(atomic.SwapInt32(&maxThreads, int32(n)))
}

// WaitThreads waits until the number of threads meet the SetMaxThreads
// settings. The function always returns true if the ctx is not canceled.
// Otherwise returns true only if the Wait is successed in the last check.
func WaitThreads(ctx context.Context) (ok bool) {
	for {
		select {
		case <-ctx.Done():
			if NumThreads() <= SetMaxThreads(0) {
				ok = true
			}
			return
		default:
			if NumThreads() > SetMaxThreads(0) {
				continue
			}
			ok = true
			return
		}
	}
}

func checkwork() error {
	_, err := exec.Command("bash", "-c", cmdThreads).Output()
	if err != nil {
		return fmt.Errorf("runtime: failed to use the package: %w", err)
	}
	return nil
}
