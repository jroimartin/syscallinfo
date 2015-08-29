// Copyright 2015 The syscallinfo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/jroimartin/syscallinfo"
)

var filename = flag.String("output", "", "output file name (standard output if omitted)")

type SyscallinfoPackage struct {
	PkgName  string
	Syscalls []syscallinfo.Syscall
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 2 {
		usage()
	}
	pkgname := flag.Arg(0)
	ctxfile := flag.Arg(1)

	ctxdata, err := ioutil.ReadFile(ctxfile)
	if err != nil {
		log.Fatalln(err)
	}

	sipkg := SyscallinfoPackage{PkgName: pkgname}
	if err := json.Unmarshal(ctxdata, &sipkg.Syscalls); err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	t := template.Must(template.New("src").Parse(srcTemplate))
	if err := t.Execute(&buf, sipkg); err != nil {
		log.Fatalln(err)
	}
	if *filename != "" {
		err = ioutil.WriteFile(*filename, buf.Bytes(), 0644)
	} else {
		_, err = os.Stdout.Write(buf.Bytes())
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run mksyscalltable.go [flags] pkgname ctxfile")
	flag.PrintDefaults()
	os.Exit(2)
}

const srcTemplate = `// MACHINE GENERATED BY 'go generate' COMMAND; DO NOT EDIT

package {{.PkgName}}

import "github.com/jroimartin/syscallinfo"

var SyscallTable = syscallinfo.SyscallTable{
{{range .Syscalls}}	{{.Num}}: syscallinfo.Syscall{
		Num: {{.Num}},
		Name: "{{.Name}}",
		Entry: "{{.Entry}}",
		Args: []syscallinfo.Argument{
{{range .Args}}			{
				RefCount: {{.RefCount}},
				Sig: "{{.Sig}}",
				Context: {{.Context}},
			},
{{end}}		},
	},
{{end}}}
`
