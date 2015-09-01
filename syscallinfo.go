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

// Initialize the package's DefaultContextHandler with default handlers.
func init() {
	Handle(CtxFD, func(n uint64) (string, error) {
		return fmt.Sprintf("%d", n), nil
	})
}

// A Syscall contains information about a syscall in a way that is OS and arch
// independent.
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

// Context gives information about how a syscall argument or return value is
// used. For instance, it specifies if an argument is a file descriptor.
type Context int

const (
	// CtxNone represents an Unknown context
	CtxNone Context = iota
	// CtxFD represents a File descriptor
	CtxFD
)

// UnmarshalJSON implements JSON unmarshaling for context.
func (ctx *Context) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("context should be a string, got %s", data)
	}
	switch s {
	case "FD":
		*ctx = CtxFD
	default:
		*ctx = CtxNone
	}
	return nil
}

// A SyscallTable contains the information about the syscalls of a specific
// OS.
type SyscallTable map[int]Syscall

// A Resolver allows to access information from a given syscall table.
type Resolver struct {
	tbl SyscallTable
}

// NewResolver returns a syscall resolver for the specified syscall table.
func NewResolver(tbl SyscallTable) Resolver {
	return Resolver{tbl: tbl}
}

// SyscallN returns a Syscall object which number matches the provided one.
func (r Resolver) SyscallN(n int) (Syscall, error) {
	sc, ok := r.tbl[n]
	if ok {
		return sc, nil
	}
	return Syscall{}, errors.New("unknown syscall")
}

// SyscallEntry returns a Syscall object which entry point matches the provided
// one.
func (r Resolver) SyscallEntry(entry string) (Syscall, error) {
	for _, sc := range r.tbl {
		if sc.Entry == entry {
			return sc, nil
		}
	}
	return Syscall{}, errors.New("unknown syscall")
}

// HandlerFunc is a function that implements how a value must be
// contextualized.
type HandlerFunc func(n uint64) (string, error)

// A ContextHandler is map that links contexts with handlers.
type ContextHandler map[Context]HandlerFunc

// Handle assigns a HandlerFunc to a context within a context handler.
func (ch ContextHandler) Handle(ctx Context, h HandlerFunc) {
	ch[ctx] = h
}

// DefaultContextHandler is the default context handler used by syscallinfo. It
// defines a basic handler for each specific context.
var DefaultContextHandler = ContextHandler{}

// Handle assigns a HandlerFunc to a context within the DefaultContextHandler.
func Handle(ctx Context, h HandlerFunc) {
	DefaultContextHandler.Handle(ctx, h)
}

// A SyscallCall represents a call to a syscall, with its own return value,
// arguments and context handler.
type SyscallCall struct {
	sc   Syscall
	ret  uint64
	args []uint64
	ch   ContextHandler
}

// NewSyscallCall returns a reference to a new SyscallCall object. The number
// of provided arguments must be greater or equal to the number of arguments
// required by the syscall.
func NewSyscallCall(sc Syscall, ret uint64, args ...uint64) (*SyscallCall, error) {
	if len(args) < len(sc.Args) {
		return nil, errors.New("invalid number of arguments")
	}
	scc := &SyscallCall{
		sc:   sc,
		args: args,
		ret:  ret,
		ch:   DefaultContextHandler,
	}
	return scc, nil
}

// SetContextHandler allows to set a custom ContextHandler to a SyscallCall
// object.
func (scc *SyscallCall) SetContextHandler(ch ContextHandler) {
	scc.ch = ch
}

// OutputOption allow to configure the output type.
type OutputOption int

const (
	// OutRet makes the ret value to be printed.
	OutRet OutputOption = 1 << iota
)

// Output returns a string with the representation of the call.
func (scc *SyscallCall) Output(opts OutputOption) (string, error) {
	str := ""
	argsStr := ""
	for i := range scc.sc.Args {
		argStr, err := scc.handleContext(scc.args[i], scc.sc.Args[i].Context)
		if err != nil {
			return "", err
		}
		argsStr += argStr + ", "
	}
	argsStr = strings.TrimSuffix(argsStr, ", ")
	str += fmt.Sprintf("%s(%s)", scc.sc.Name, argsStr)
	if opts&OutRet != 0 {
		retStr, err := scc.handleContext(scc.ret, scc.sc.Context)
		if err != nil {
			return "", err
		}
		str += " = " + retStr
	}
	return str, nil
}

// String returns a string with the representation of the call plus the return
// value. An empty string is returned on error.
func (scc *SyscallCall) String() string {
	str, err := scc.Output(OutRet)
	if err != nil {
		return ""
	}
	return str
}

// handleContext returns a string with the contextualized representation of the
// provided value.
func (scc *SyscallCall) handleContext(n uint64, ctx Context) (string, error) {
	h, ok := scc.ch[ctx]
	if ok && h != nil {
		return h(n)
	}

	// Fallback to DefaultContextHandler
	h, ok = DefaultContextHandler[ctx]
	if ok && h != nil {
		return h(n)
	}

	// Default value representation
	return fmt.Sprintf("%#08x", n), nil
}
