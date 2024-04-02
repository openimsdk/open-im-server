package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
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

	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Randomly select two ports between 10000 and 20000
	port1 := rand.Intn(10001) + 10000 // This generates a number between 0-10000, then adds 10000
	port2 := rand.Intn(10001) + 10000

	// Ensure port2 is different from port1
	for port2 == port1 {
		port2 = rand.Intn(10001) + 10000
	}

	// Print the selected ports
	fmt.Printf("Randomly selected TCP ports to listen on: %d, %d\n", port1, port2)

	// Function to start a TCP listener on a specified port
	startListener := func(port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Printf("Error starting TCP listener on port %d: %v\n", port, err)
			return
		}
		defer listener.Close()
		fmt.Printf("Listening on port %d...\n", port)

		// Here we simply keep the listener running. In a real application, you would accept connections.
		select {} // Block forever
	}

	// Start TCP listeners on the selected ports
	go startListener(port1)
	go startListener(port2)

	// Block forever
	select {}
}
