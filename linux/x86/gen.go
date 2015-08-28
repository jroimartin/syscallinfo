package x86

//go:generate go run $GOPATH/src/github.com/jroimartin/syscallinfo/mksyscallinfo.go -arch x86 -output syscallinfo.go syscall_32.json
