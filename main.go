package main

import (
	"fmt"
	"net"
)

const (
	DiscoveryPort = "5501"
	DiscoveryMsg  = "P2P_CHAT_DISCOVERY"
)

func main() {
	fmt.Println("LOCAL P2PChat: Discovering clients...")
	listenAddr, err := net.ResolveUDPAddr("udp", ":"+DiscoveryPort)
	if err != nil {
		fmt.Printf("Failed to resolve listen address: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+DiscoveryPort)
	if err != nil {
		panic(err)
	}

	go func() {
		_, err := conn.WriteToUDP([]byte(DiscoveryMsg), broadcastAddr)
		if err != nil {
			fmt.Printf("Broadcast err: %v\n", err)
		}
	}()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("❌ Error reading packet: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		if message == DiscoveryMsg {
			fmt.Printf("[PEER DISCOVERED] %s\n", addr.String())
		}
	}
}
