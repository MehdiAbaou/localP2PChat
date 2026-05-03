package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {

	fmt.Print("Enter your name (Max 10 chars): ")
	var name string
	fmt.Scanln(&name)

	if len(name) > 10 {
		name = name[:10]
	}

	// Listen on port 5501
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 5501})
	defer conn.Close()

	target, _ := net.ResolveUDPAddr("udp", "255.255.255.255:5501")

	var peers []string

	go func() {
		conn.WriteToUDP([]byte("P2PDSCVMSG/"+name+"/"+conn.LocalAddr().String()), target)
		fmt.Println("Discovering...")

		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			if len(strings.Split(string(buf[:n]), "/")) < 3 {
				continue
			}
			if strings.Split(string(buf[:n]), "/")[0] == "P2PDSCVMSG" && strings.Split(string(buf[:n]), "/")[1] != name {
				fmt.Printf("Found Peer: %s (%s)\n", strings.Split(string(buf[:n]), "/")[1], strings.Split(string(buf[:n]), "/")[2])
				peers = append(peers, strings.Split(string(buf[:n]), "/")[1])
			}
		}
	}()

	for {
		var input string
		fmt.Scanln(&input)
		conn.WriteToUDP([]byte(input), target)
	}
}
