package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	// Initialize the random seed
	rand.Seed(time.Now().UnixNano())

	// Randomly select two ports between 10000 and 20000
	port1 := rand.Intn(10001) + 10000
	port2 := rand.Intn(10001) + 10000
	// Ensure port2 is different from port1
	for port2 == port1 {
		port2 = rand.Intn(10001) + 10000
	}

	fmt.Printf("Randomly selected TCP ports to listen on: %d, %d\n", port1, port2)

	// Listen on the selected ports
	go listenOnPort(port1)
	go listenOnPort(port2)

	// Wait for a termination signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
}

func listenOnPort(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error listening on port %d: %v", port, err)
	}
	defer l.Close()

	fmt.Printf("Listening on port %d...\n", port)

	// Normally, you would accept connections with l.Accept() in a loop here.
	// For this example, just sleep indefinitely to simulate a service.
	select {}

}
