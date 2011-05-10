// Usage: go [go-file]
// This is useful while developing & testing 
// a Go program. A single command that'll
// take care of the compiling and running
// of the specified program. If the source 
// file is newer than the binary, it'll compile 
// it and run it. Otherwise, it'll just run it.
// If the Go program is FOO.go, the binary
// is named FOO.out.

package main

import (
	"os"
	"flag"
	"fmt"
	"strings"
	"exec"
)

func fatalError(msg string) {
	if msg != "" {
		fmt.Printf("%s\n", msg)
	}
	os.Exit(1)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Returns true if A is newer that B.
// File A is assumed to exist.
// If File B doesn't exist, return true.
// Otherwise their last modified times are
// compared.
func isFileNewer(a, b string) bool {
	as, e1 := os.Stat(a)
	bs, e2 := os.Stat(b)
	if e1 != nil {
		fatalError(e1.String())
	}
	if e2 != nil {
		return true
	}
	return as.Mtime_ns > bs.Mtime_ns
}

func run(prog string, argv []string) {
	fullp, err := exec.LookPath(prog)
	argvP := make([]string, len(argv) + 1)
	argvP[0] = fullp
	for i := 0; i < len(argv); i++ {
		argvP[i + 1] = argv[i]
	}
	if err != nil {
		fatalError(err.String())
	}
	cmd, err := exec.Run(
		fullp, argvP, nil, "", exec.DevNull, 
		exec.PassThrough, exec.MergeWithStdout)
	if err != nil {
		fatalError(err.String())
	}
	err = cmd.Close()
	if err != nil {
		fatalError(err.String())
	}
}

func filenameNoExt(name string) string {
	s := strings.Split(name, ".", -1)
	t := s[0:len(s)-1]
	return strings.Join(t, ".")
}

func main() {
	if flag.NArg() != 1 {
		fatalError("usage: go [Go-program]")
	}

	goProg := flag.Arg(0)
	if hasSuf := strings.HasSuffix(goProg, ".go"); !hasSuf {
		goProg += ".go"
	}

	if exists := fileExists(goProg); !exists {
		fatalError(goProg + " doesn't exist")
	}

	exeName := goProg + ".exe"
	lnName := filenameNoExt(goProg) + ".8"
	needsCompile := isFileNewer(goProg, exeName)
	
	if needsCompile {
		run("8g", []string{goProg})
		run("8l", []string{lnName})
		os.Rename("8.out", exeName)
		os.Remove(lnName) 
	}

	//oArgs := flag.Args()
	//pArgs := oArgs[1:len(oArgs)]
	//fmt.Printf("%+v\n", pArgs)
	run("./" + exeName, nil)
}
