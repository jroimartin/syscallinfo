package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var filename = flag.String("output", "", "output file name (standard output if omitted)")

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

type Argument struct {
	RefCount int
	Sig      string
	Context  Context
}

type Syscall struct {
	Num   int
	Name  string
	Entry string
	Args  []Argument
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "no file to parse provided\n")
		usage()
	}
	ctxfile := flag.Arg(0)

	ctxdata, err := ioutil.ReadFile(ctxfile)
	if err != nil {
		log.Fatalln(err)
	}

	var syscalls []Syscall
	if err := json.Unmarshal(ctxdata, &syscalls); err != nil {
		log.Fatalln(err)
	}

	data := generateOutput(syscalls)

	if *filename != "" {
		err = ioutil.WriteFile(*filename, data, 0644)
	} else {
		_, err = os.Stdout.Write(data)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func generateOutput(syscalls []Syscall) []byte {
	// TODO(jrm): Print code
	var buf bytes.Buffer
	for _, s := range syscalls {
		fmt.Fprintf(&buf, "%#v\n", s)
	}
	return buf.Bytes()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mksyscallinfo_linux [flags] ctxfile\n")
	flag.PrintDefaults()
	os.Exit(2)
}
