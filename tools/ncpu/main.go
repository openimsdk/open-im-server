package main

import (
	"fmt"
	"runtime"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	maxprocs.Set()
	fmt.Print(runtime.GOMAXPROCS(0))
}
