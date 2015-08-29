// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amd64

//go:generate go run $GOPATH/src/github.com/jroimartin/syscallinfo/mksyscallinfo.go -arch amd64 -output syscallinfo.go syscall_64.json
