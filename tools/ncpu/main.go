package main

import (
	"fmt"
	"runtime"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// Set maxprocs with a custom logger that does nothing to ignore logs.
	maxprocs.Set(maxprocs.Logger(func(string, ...interface{}) {
		// Intentionally left blank to suppress all log output from automaxprocs.
	}))

	// Now this will print the GOMAXPROCS value without printing the automaxprocs log message.
	fmt.Println(runtime.GOMAXPROCS(0))
}
