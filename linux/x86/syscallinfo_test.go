package x86

import "testing"

var checks = []struct {
	num      int
	entry    string
	nargs    int
	nilError bool
}{
	{26, "sys_ptrace", 4, true},
	{666, "", 0, false},
}

func TestGetSyscall(t *testing.T) {
	for _, c := range checks {
		sc, err := GetSyscall(c.num)
		if err != nil {
			if c.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if sc.Entry != c.entry {
			t.Errorf("wrong entry (want=%v, get=%v)", c.entry, sc.Entry)
		}
		if len(sc.Args) != c.nargs {
			t.Errorf("wrong number of arguments (want=%v, get=%v)", c.nargs, len(sc.Args))
		}
	}
}

func TestGetSyscallByEntry(t *testing.T) {
	for _, c := range checks {
		sc, err := GetSyscallByEntry(c.entry)
		if err != nil {
			if c.nilError {
				t.Errorf("wrong error (want=nil, get=%v)", err)
			}
			continue
		}
		if sc.Num != c.num {
			t.Errorf("wrong number (want=%v, get=%v)", c.num, sc.Num)
		}
		if len(sc.Args) != c.nargs {
			t.Errorf("wrong number of arguments (want=%v, get=%v)", c.nargs, len(sc.Args))
		}
	}
}
