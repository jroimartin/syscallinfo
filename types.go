// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo

import (
	"encoding/json"
	"fmt"
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
