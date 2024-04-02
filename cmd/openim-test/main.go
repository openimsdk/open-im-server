package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// (Your flag definitions and random port selection logic here)

	// Initialize a channel to listen for interrupt signals
	sigs := make(chan os.Signal, 1)
	// Register the channel to receive SIGINT and SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Function to start a TCP listener on a specified port
	startListener := func(port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Printf("Error starting TCP listener on port %d: %v\n", port, err)
			return
		}
		defer listener.Close()
		fmt.Printf("Listening on port %d...\n", port)

		// Accept connections in a loop
		go func() {
			for {
				_, err := listener.Accept()
				if err != nil {
					fmt.Printf("Error accepting connection on port %d: %v\n", port, err)
					break
				}
				// Handle connection (e.g., close immediately for this example)
				// conn.Close()
			}
		}()
	}

	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Randomly select two ports between 10000 and 20000
	port1 := rand.Intn(10001) + 10000 // This generates a number between 0-10000, then adds 10000
	port2 := rand.Intn(10001) + 10000

	// Ensure port2 is different from port1
	for port2 == port1 {
		port2 = rand.Intn(10001) + 10000
	}

	// Start TCP listeners on the selected ports
	// (Assuming port1 and port2 are defined as per your logic)
	go startListener(port1)
	go startListener(port2)

	// Wait for an interrupt signal
	<-sigs
	fmt.Println("Shutting down...")
}
