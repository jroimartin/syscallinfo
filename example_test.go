// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo_test

import (
	"fmt"

	"github.com/jroimartin/syscallinfo"
	"github.com/jroimartin/syscallinfo/linux_386"
)

func ExampleSyscall_Repr() {
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

func ExampleResolver_Handle() {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	r.Handle(syscallinfo.CTX_FD, func(n uint64) (string, error) {
		return fmt.Sprintf("FD(%d)", n), nil
	})
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
	// read(FD(1), 0x00000002, 0x00000003) = 0x00000004
}
