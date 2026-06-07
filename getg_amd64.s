// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

//go:build goidfast

#include "textflag.h"

// func getg() unsafe.Pointer
//
// On amd64 the goroutine pointer lives in thread-local storage. go_tls.h
// is not includable from non-runtime packages, so this inlines what its
// get_tls(CX) / g(CX) macros expand to: load the TLS base, then read g at
// 0(CX)(TLS*1). The assembler lowers the TLS pseudo-register per OS.
TEXT ·getg(SB), NOSPLIT, $0-8
	MOVQ TLS, CX
	MOVQ 0(CX)(TLS*1), AX
	MOVQ AX, ret+0(FP)
	RET
