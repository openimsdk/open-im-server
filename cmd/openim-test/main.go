package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define flags
	index := flag.Int("i", 0, "Index number")
	config := flag.String("c", "", "Configuration path")

	// Parse the flags
	flag.Parse()

	// Print the values of the flags
	fmt.Printf("-i %d -c %s\n", *index, *config)
}
