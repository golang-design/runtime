// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build goidfast

#include "textflag.h"

// func getg() unsafe.Pointer
// On arm64 the goroutine pointer is held in register R28 (the g register).
TEXT ·getg(SB), NOSPLIT, $0-8
	MOVD g, R0
	MOVD R0, ret+0(FP)
	RET
