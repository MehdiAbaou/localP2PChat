package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on port 5501
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 5501})
	defer conn.Close()

	target, _ := net.ResolveUDPAddr("udp", "255.255.255.255:5501")

	go func() {
		conn.WriteToUDP([]byte("I_AM_HERE"), target)
	}()

	fmt.Println("Discovering...")

	buf := make([]byte, 1024)
	for {
		n, addr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("Found Peer: %s (%s)\n", addr, string(buf[:n]))
	}
}
