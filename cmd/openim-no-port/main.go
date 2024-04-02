package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// Define flags
	index := flag.Int("i", 0, "Index number")
	config := flag.String("c", "", "Configuration path")

	// Parse the flags
	flag.Parse()
	// Print the values of the flags
	fmt.Printf("args: -i %d -c %s\n", *index, *config)
	time.Sleep(500 * time.Second)
}
