package logger

import (
	"fmt"
	"runtime"
)

// var pcs []uintptr

func init() {
	// pcs = make([]uintptr, 20)
	fmt.Println("init func")
}

func caller(skip int) (string, string, int, bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", -1, false
	}
	f := runtime.FuncForPC(pc)
	return file, f.Name(), line, ok
}
