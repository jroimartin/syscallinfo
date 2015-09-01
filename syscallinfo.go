// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// A Syscall contains information about a syscall in a way that it
// is OS and arch independent.
type Syscall struct {
	// Num is the number of the syscall.
	Num int

	// Name is the user-readable name of the syscall.
	Name string

	// Entry is the entry point of the syscall (function name).
	Entry string

	// Context is specifies under which context the syscall's return value is
	// used.
	Context Context

	// Args is a slice containing all the syscall's argurments.
	Args []Argument

	// r is the resolver used to resolve this syscall.
	r Resolver
}

// Argument represents a syscall argument.
type Argument struct {
	// RefCount is the level of indirection.
	RefCount int

	// Sig is the signature of the argument (type and name).
	Sig string

	// Context specifies under which context this argument is used.
	Context Context
}

// Context gives information about how a syscall argument is used. For
// instance, it specifies if the argument is a file descriptor.
type Context int

const (
	CTX_NONE Context = iota // Unknown context
	CTX_FD                  // File descriptor
)

// UnmarshalJSON implements JSON unmarshaling for context.
func (ctx *Context) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("context should be a string, got %s", data)
	}
	switch s {
	case "FD":
		*ctx = CTX_FD
	default:
		*ctx = CTX_NONE
	}
	return nil
}

// A SyscallTable contains the information about the syscalls of a specific
// OS.
type SyscallTable map[int]Syscall

// A Resolver allows to access information from a given syscall table.
type Resolver struct {
	tbl          SyscallTable
	handlerFuncs map[Context]HandlerFunc
}

// HandlerFunc is a function that implements how a value must be
// contextualized.
type HandlerFunc func(n uint64) (string, error)

// NewResolver returns a syscall resolver for the specified syscall table.
func NewResolver(tbl SyscallTable) Resolver {
	return Resolver{
		tbl:          tbl,
		handlerFuncs: make(map[Context]HandlerFunc),
	}
}

// Syscall returns a Syscall object which number matches the provided one.
func (r Resolver) Syscall(n int) (Syscall, error) {
	sc, ok := r.tbl[n]
	if ok {
		sc.r = r
		return sc, nil
	}
	return Syscall{}, errors.New("unknown syscall")
}

// SyscallByEntry returns a Syscall object which entry point matches the
// provided one.
func (r Resolver) SyscallByEntry(entry string) (Syscall, error) {
	for _, sc := range r.tbl {
		if sc.Entry == entry {
			sc.r = r
			return sc, nil
		}
	}
	return Syscall{}, errors.New("unknown syscall")
}

// Handle assigns a HandlerFunc to a context.
func (r Resolver) Handle(ctx Context, h HandlerFunc) {
	r.handlerFuncs[ctx] = h
}

// handleContext returns a string with the contextualized representation of the
// provided value.
func (r Resolver) handleContext(n uint64, ctx Context) (string, error) {
	h, ok := r.handlerFuncs[ctx]
	if ok && h != nil {
		return h(n)
	}
	switch ctx {
	case CTX_FD:
		return fmt.Sprintf("%d", n), nil
	default:
		return fmt.Sprintf("%#08x", n), nil
	}
}

// Repr returns a string with the representation of the call plus the return
// value. The number of provided arguments must be greater or equal to the
// number of arguments required by the syscall.
func (sc Syscall) Repr(ret uint64, args ...uint64) (string, error) {
	callStr, err := sc.ReprCall(args...)
	if err != nil {
		return "", err
	}
	retStr, err := sc.r.handleContext(ret, sc.Context)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %s", callStr, retStr), nil
}

// ReprCall returns a string with the representation of the call. The number of
// provided arguments must be greater or equal to the number of arguments
// required by the syscall.
func (sc Syscall) ReprCall(args ...uint64) (string, error) {
	if len(args) < len(sc.Args) {
		return "", errors.New("invalid number of arguments")
	}
	argsStr := ""
	for i := range sc.Args {
		argStr, err := sc.r.handleContext(args[i], sc.Args[i].Context)
		if err != nil {
			return "", err
		}
		argsStr += argStr + ", "
	}
	argsStr = strings.TrimSuffix(argsStr, ", ")
	return fmt.Sprintf("%s(%s)", sc.Name, argsStr), nil
}
