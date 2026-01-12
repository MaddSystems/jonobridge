package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on UDP port 12345
	conn, err := net.ListenPacket("udp", ":8690")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		// Read incoming UDP message
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Print sender IP and received message
		fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), addr.String())
	}
}
