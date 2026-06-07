// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goidfast && (amd64 || arm64)

package runtime

import (
	"sync"
	"testing"
)

// TestGoidCalibrated guards the fast path: on a supported arch the goidfast
// build must successfully calibrate, otherwise it silently degrades to the
// slow path and the tag buys nothing.
func TestGoidCalibrated(t *testing.T) {
	if !goidCalibrated {
		t.Fatal("goidfast build failed to calibrate the goid offset")
	}
}

// TestGoidFastMatchesSlow verifies the fast id equals the authoritative
// slow id, on the current goroutine and across many fresh ones.
func TestGoidFastMatchesSlow(t *testing.T) {
	if got, want := Goid(), goidSlow(); got != want {
		t.Fatalf("fast Goid = %d, slow = %d", got, want)
	}

	var wg sync.WaitGroup
	errs := make(chan [2]uint64, 64)
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if fast, slow := Goid(), goidSlow(); fast != slow {
				errs <- [2]uint64{fast, slow}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Fatalf("fast Goid = %d, slow = %d in goroutine", e[0], e[1])
	}
}
