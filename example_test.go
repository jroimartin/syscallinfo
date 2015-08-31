// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo_test

import (
	"fmt"

	"github.com/jroimartin/syscallinfo"
	"github.com/jroimartin/syscallinfo/linux_386"
)

func ExampleRepr() {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	sc, err := r.Syscall(3)
	if err != nil {
		return
	}
	str, err := sc.Repr(4, 1, 2, 3)
	if err != nil {
		return
	}
	fmt.Println(str)

	// Output:
	// read(1, 0x00000002, 0x00000003) = 0x00000004
}
