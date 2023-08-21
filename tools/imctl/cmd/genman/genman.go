package main

import (
	"fmt"
	"os"

	"k8s.io/kubernetes/cmd/genutils"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/man/man1"
	module := ""
	if len(os.Args) == 3 {
		path = os.Args[1]
		module = os.Args[2]
	} else {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory] [module] \n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by command so the output is consistent,
	// regardless of where we run.
	os.Setenv("HOME", "/home/username")
}
