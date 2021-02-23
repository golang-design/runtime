// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.design/x/runtime"
)

func TestSetMaxThreads(t *testing.T) {
	runtime.SetMaxThreads(10)

	// create a lot of threads by sleep gs
	wg := sync.WaitGroup{}
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go func() {
			time.Sleep(time.Second * 1)
			wg.Done()
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ok := runtime.WaitThreads(ctx)
	if !ok {
		t.Fatal("mkill failed in 10s")
	}
	wg.Wait()
}

func TestMinThreads(t *testing.T) {
	old := runtime.SetMaxThreads(0)
	n := runtime.NumCPU()
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
