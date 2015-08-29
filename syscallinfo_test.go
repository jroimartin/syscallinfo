// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo_test

import (
	"testing"

	"github.com/jroimartin/syscallinfo"
	"github.com/jroimartin/syscallinfo/linux32"
)

var checksResolution = []struct {
	num      int
	entry    string
	nargs    int
	context  []syscallinfo.Context
	nilError bool
}{
	{
		3,
		"sys_read",
		3,
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
		0,
		[]syscallinfo.Context{},
		false,
	},
}

func TestSyscall(t *testing.T) {
	r := syscallinfo.NewResolver(linux32.SyscallTable)
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
		if len(sc.Args) != check.nargs {
			t.Errorf("wrong number of arguments (want=%v, get=%v)", check.nargs, len(sc.Args))
			continue
		}
		for i := range sc.Args {
			if sc.Args[i].Context != check.context[i] {
				t.Errorf("wrong context (want=%v, get=%v)", check.context[i], sc.Args[i].Context)
			}
		}
	}
}

func TestSyscallByEntry(t *testing.T) {
	r := syscallinfo.NewResolver(linux32.SyscallTable)
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
		if len(sc.Args) != check.nargs {
			t.Errorf("wrong number of arguments (want=%v, get=%v)", check.nargs, len(sc.Args))
			continue
		}
		for i := range sc.Args {
			if sc.Args[i].Context != check.context[i] {
				t.Errorf("wrong context (want=%v, get=%v)", check.context[i], sc.Args[i].Context)
			}
		}
	}
}

var checksReprs = []struct {
	num      int
	args     []uint64
	want     string
	nilError bool
}{
	{
		3,
		[]uint64{1, 2, 3},
		"read(1, 0x00000002, 0x00000003)",
		true,
	},
	{
		3,
		[]uint64{1, 2},
		"",
		false,
	},
}

func TestRepr(t *testing.T) {
	r := syscallinfo.NewResolver(linux32.SyscallTable)

	for _, check := range checksReprs {
		str, err := r.Repr(check.num, check.args...)
		if err != nil {
			if check.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if str != check.want {
			t.Errorf("wrong string (want=%v, get=%v)", check.want, str)
		}
	}
}
