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
	"sync"
	"sync/atomic"
	"time"
)

var (
	pid = os.Getpid()
	// minimum number of threads required by the runtime
	minThreads = int32(NumCPU()) + 2
	// 2 meaning runtime sysmon thread + template thread
	maxThreads = int32(NumCPU()) + 2
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
	// Start the watcher exactly once, on the first call. The previous
	// code inverted this condition (it ran the start block on every call
	// except the first), so a lone SetMaxThreads call never started the
	// killer and every subsequent call leaked another ticker goroutine.
	if atomic.CompareAndSwapInt32(&started, 0, 1) {
		if err := checkwork(); err != nil {
			return 0
		}
		if debug {
			fmt.Printf("runtime: pid %v, maxThread %v, interval %v\n",
				pid, maxThreads, interval)
		}
		go watch()
	}

	if n < int(minThreads) {
		return int(atomic.LoadInt32(&maxThreads))
	}

	return int(atomic.SwapInt32(&maxThreads, int32(n)))
}

// watch periodically kills surplus OS threads so their number stays at or
// below the configured maximum. It runs for the lifetime of the program.
func watch() {
	t := time.NewTicker(interval)
	for range t.C {
		n := NumThreads()
		nkill := int32(n) - atomic.LoadInt32(&maxThreads)
		if nkill <= 0 {
			if debug {
				fmt.Printf("runtime: checked #threads total %v / max %v\n",
					n, atomic.LoadInt32(&maxThreads))
			}
			continue
		}
		// A goroutine that exits while locked to its OS thread takes the
		// thread down with it.
		wg := sync.WaitGroup{}
		wg.Add(int(nkill))
		for i := int32(0); i < nkill; i++ {
			go func() {
				LockOSThread()
				wg.Done()
			}()
		}
		wg.Wait()
		if debug {
			fmt.Printf("runtime: killing #threads, remaining: %v\n", n)
		}
	}
}

// WaitThreads waits until the number of threads meet the SetMaxThreads
// settings. The function always returns true if the ctx is not canceled.
// Otherwise returns true only if the Wait is successed in the last check.
func WaitThreads(ctx context.Context) (ok bool) {
	// Poll on a bounded interval instead of busy-spinning. The previous
	// loop spun with no backoff and called SetMaxThreads(0) as a getter,
	// forking a shell (and, given the start-up bug above, leaking a
	// watcher) on every iteration.
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if NumThreads() <= int(atomic.LoadInt32(&maxThreads)) {
			return true
		}
		select {
		case <-ctx.Done():
			return NumThreads() <= int(atomic.LoadInt32(&maxThreads))
		case <-t.C:
		}
	}
}

// checkwork verifies, before the watcher starts, that this platform can
// actually count the process's threads (numThreads is implemented per
// platform in mkill_<goos>.go).
func checkwork() error {
	if numThreads() <= 0 {
		return fmt.Errorf("runtime: cannot determine the number of OS threads on this platform")
	}
	return nil
}
