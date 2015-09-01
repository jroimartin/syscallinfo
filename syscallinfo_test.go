// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo_test

import (
	"fmt"
	"testing"

	"github.com/jroimartin/syscallinfo"
	"github.com/jroimartin/syscallinfo/linux_386"
)

var checksResolution = []struct {
	num      int
	entry    string
	context  []syscallinfo.Context
	nilError bool
}{
	{
		3,
		"sys_read",
		[]syscallinfo.Context{
			syscallinfo.CTX_FD,
			syscallinfo.CTX_NONE,
			syscallinfo.CTX_NONE,
		},
		true,
	},
	{
		666,
		"",
		[]syscallinfo.Context{},
		false,
	},
}

func TestSyscall(t *testing.T) {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	for _, check := range checksResolution {
		sc, err := r.Syscall(check.num)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if sc.Entry != check.entry {
			t.Errorf("wrong entry (want=%v, get=%v)", check.entry, sc.Entry)
		}
		if len(sc.Args) != len(check.context) {
			t.Errorf("wrong number of arguments (want=%v, get=%v)",
				len(check.context), len(sc.Args))
			continue
		}
		for i := range sc.Args {
			if sc.Args[i].Context != check.context[i] {
				t.Errorf("wrong context (want=%v, get=%v)",
					check.context[i], sc.Args[i].Context)
			}
		}
	}
}

func TestSyscallByEntry(t *testing.T) {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	for _, check := range checksResolution {
		sc, err := r.SyscallByEntry(check.entry)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if sc.Num != check.num {
			t.Errorf("wrong number (want=%v, get=%v)", check.num, sc.Num)
		}
		if len(sc.Args) != len(check.context) {
			t.Errorf("wrong number of arguments (want=%v, get=%v)",
				len(check.context), len(sc.Args))
			continue
		}
		for i := range sc.Args {
			if sc.Args[i].Context != check.context[i] {
				t.Errorf("wrong context (want=%v, get=%v)",
					check.context[i], sc.Args[i].Context)
			}
		}
	}
}

var checksReprs = []struct {
	num      int
	args     []uint64
	retval   uint64
	reprcall string
	repr     string
	nilError bool
}{
	{
		3,
		[]uint64{1, 2, 3},
		4,
		"read(1, 0x00000002, 0x00000003)",
		"read(1, 0x00000002, 0x00000003) = 0x00000004",
		true,
	},
	{
		3,
		[]uint64{1, 2},
		3,
		"",
		"",
		false,
	},
	{
		5,
		[]uint64{1, 2, 3},
		4,
		"open(0x00000001, 0x00000002, 0x00000003)",
		"open(0x00000001, 0x00000002, 0x00000003) = 4",
		true,
	},
}

func TestReprCall(t *testing.T) {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)

	for _, check := range checksReprs {
		sc, err := r.Syscall(check.num)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		str, err := sc.ReprCall(check.args...)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if str != check.reprcall {
			t.Errorf("wrong string (want=%v, get=%v)", check.reprcall, str)
		}
	}
}

func TestRepr(t *testing.T) {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)

	for _, check := range checksReprs {
		sc, err := r.Syscall(check.num)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		str, err := sc.Repr(check.retval, check.args...)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if str != check.repr {
			t.Errorf("wrong string (want=%v, get=%v)", check.repr, str)
		}
	}
}

var checkHandle = struct {
	num      int
	args     []uint64
	retval   uint64
	reprcall string
}{
	3,
	[]uint64{1, 2, 3},
	4,
	"read(test-1, 0x00000002, 0x00000003) = 0x00000004",
}

func TestHandle(t *testing.T) {
	r := syscallinfo.NewResolver(linux_386.SyscallTable)
	r.Handle(syscallinfo.CTX_FD, func(n uint64) (string, error) {
		return fmt.Sprintf("test-%d", n), nil
	})

	sc, err := r.Syscall(checkHandle.num)
	if err != nil {
		t.Errorf("wrong error (want=nil, get=%v)", err)
		return
	}
	str, err := sc.Repr(checkHandle.retval, checkHandle.args...)
	if err != nil {
		t.Errorf("wrong error (want=nil, get=%v)", err)
		return
	}
	if str != checkHandle.reprcall {
		t.Errorf("wrong string (want=%v, get=%v)", checkHandle.reprcall, str)
	}
}
