// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscallinfo

import (
	"encoding/json"
	"fmt"
)

type Syscall struct {
	Num   int
	Name  string
	Entry string
	Args  []Argument
}

type Argument struct {
	RefCount int
	Sig      string
	Context  Context
}

type Context int

const (
	CTX_NONE Context = iota
	CTX_FD
)

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
