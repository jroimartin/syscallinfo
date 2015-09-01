// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo_test

import (
	"fmt"

	"github.com/jroimartin/syscallinfo"
	"github.com/jroimartin/syscallinfo/linux_386"
)

func ExampleNewSyscallCall() {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	sc, err := r.SyscallN(3)
	if err != nil {
		return
	}
	scc, err := syscallinfo.NewSyscallCall(sc, 4, 1, 2, 3)
	if err != nil {
		return
	}
	fmt.Println(scc)

	// Output:
	// read(1, 0x00000002, 0x00000003) = 0x00000004
}

func ExampleHandle() {
	syscallinfo.Handle(syscallinfo.CtxFD, func(n uint64) (string, error) {
		return fmt.Sprintf("FD(%d)", n), nil
	})

	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	sc, err := r.SyscallN(3)
	if err != nil {
		return
	}
	scc, err := syscallinfo.NewSyscallCall(sc, 4, 1, 2, 3)
	if err != nil {
		return
	}
	fmt.Println(scc)

	// Output:
	// read(FD(1), 0x00000002, 0x00000003) = 0x00000004
}
