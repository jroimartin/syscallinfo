// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

//go:generate go run $GOPATH/src/github.com/jroimartin/syscallinfo/mksyscallinfo.go -arch x86 -output syscallinfo.go syscall_32.json
