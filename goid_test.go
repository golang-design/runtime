// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package runtime_test

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"golang.design/x/runtime"
)

func slowGetGoID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseInt(idField, 10, 64) // very unlikely to be failed
	return id
}

func ExampleGoid() {
	cnums := make(chan int, 100)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			cnums <- int(runtime.Goid()) // down cast, wrong in large goid
			wg.Done()
		}()
	}
	wg.Wait()
	close(cnums)

	nums := []int{int(runtime.Goid())}
	for v := range cnums {
		nums = append(nums, v)
	}
	sort.Ints(nums)
	fmt.Printf("%v", nums)
}

func TestGoid(t *testing.T) {
	got := runtime.Goid()
	want := slowGetGoID()
	if got != uint64(want) {
		t.Errorf("want %d, got: %d", want, got)
	}
}

func BenchmarkGoid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = runtime.Goid()
	}
}
