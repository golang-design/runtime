// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime_test

import (
	"context"
	"fmt"
	stdruntime "runtime"
	"sync"
	"testing"
	"time"

	"golang.design/x/runtime"
)

func TestNumThreads(t *testing.T) {
	// This is intentionally the first test in the file: it runs before any
	// SetMaxThreads call, so the killer watcher is not yet started and the
	// count is not racing the reaper. It guards numThreads against silently
	// returning 0 (e.g. the old bash/ps shell-out failing in a container)
	// and checks that the count actually reflects live OS threads.
	base := runtime.NumThreads()
	if base <= 0 {
		t.Fatalf("NumThreads returned %d, want > 0", base)
	}

	const k = 16
	var ready sync.WaitGroup
	ready.Add(k)
	release := make(chan struct{})
	done := make(chan struct{})
	for i := 0; i < k; i++ {
		go func() {
			stdruntime.LockOSThread()
			ready.Done()
			<-release // hold the OS thread so it stays counted
			done <- struct{}{}
		}()
	}
	ready.Wait()

	got := runtime.NumThreads()
	close(release)
	for i := 0; i < k; i++ {
		<-done
	}

	// k goroutines locked to distinct OS threads require at least k threads.
	if got < k {
		t.Fatalf("NumThreads did not reflect %d locked OS threads: got %d", k, got)
	}
}

func TestSetMaxThreads(t *testing.T) {
	runtime.SetMaxThreads(10)

	// create a lot of threads by sleep g
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			stdruntime.LockOSThread()
			defer stdruntime.UnlockOSThread()
			time.Sleep(time.Second * 1)
			wg.Done()
		}()
	}
	t.Logf("current threads: %d", runtime.NumThreads())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ok := runtime.WaitThreads(ctx)
	if !ok {
		t.Fatal("mkill failed in 10s")
	}
	wg.Wait()
	t.Logf("current threads: %d", runtime.NumThreads())
}

func TestMinThreads(t *testing.T) {
	old := runtime.SetMaxThreads(0)
	n := stdruntime.NumCPU()
	if runtime.SetMaxThreads(n-1) != old {
		t.Fatalf("number of threads is less than required in the runtime")
	}
}

func ExampleSetMaxThreads() {
	runtime.SetMaxThreads(10)
	// Output:
}

func ExampleNumThreads() {
	runtime.SetMaxThreads(10)
	runtime.WaitThreads(context.Background())
	fmt.Println(runtime.NumThreads() <= 10)
	// Output:
	// true
}

func ExampleWaitThreads() {
	runtime.SetMaxThreads(10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	fmt.Println(runtime.WaitThreads(ctx))
	// Output:
	// true
}
