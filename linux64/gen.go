// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux64

//go:generate go run $GOPATH/src/github.com/jroimartin/syscallinfo/mksyscalltable.go -output syscalltable.go linux64 syscall_64.json
