package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Initialize the random seed
	rand.Seed(time.Now().UnixNano())

	// Define flags
	index := flag.Int("i", 0, "Index number")
	config := flag.String("c", "", "Configuration path")

	// Parse the flags
	flag.Parse()

	// Print the values of the flags
	fmt.Printf("args: -i %d -c %s\n", *index, *config)

	// Generate a random number (0 or 1) and subtract 1 to get 0 or -1
	randomValue := rand.Intn(2) - 1

	if randomValue == 0 {
		fmt.Println("Sleeping for 500 seconds...")
		time.Sleep(500 * time.Second)
	} else if randomValue == -1 {
		fmt.Fprintln(os.Stderr, "An error occurred. Exiting with status -1.")
		os.Exit(-1)
	}
}
