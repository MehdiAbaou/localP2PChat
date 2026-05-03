package main

import (
	"crypto/md5"
	"fmt"
	"net"
	"strings"
)

type Peer struct {
	name   string
	ipHash string
}

type Backend struct {
	conn     *net.UDPConn
	target   *net.UDPAddr
	username string
	ipHash   string
	Peers    map[string]Peer
	Messages []string
}

func NewBackend(username string) (*Backend, error) {
	hash := getIPHash()
	if hash == "unknown" {
		return nil, fmt.Errorf("could not get IP hash")
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 5501})
	if err != nil {
		return nil, err
	}

	target, err := net.ResolveUDPAddr("udp", "255.255.255.255:5501")
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Backend{
		conn:     conn,
		target:   target,
		username: username,
		ipHash:   hash,
		Peers:    make(map[string]Peer),
	}, nil
}

func (b *Backend) StartDiscovery() {
	b.conn.WriteToUDP([]byte("P2PDSCVMSG/"+b.username+"/"+b.ipHash), b.target)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := b.conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			msg := string(buf[:n])
			parts := strings.Split(msg, "/") // magicnumber, name, hash
			if len(parts) < 3 {
				continue
			}

			switch parts[0] {
			case "P2PDSCVMSG":
				b.Peers[parts[2]] = Peer{
					name:   parts[1],
					ipHash: parts[2],
				}
			case "P2PLEAVMSG":
				delete(b.Peers, parts[2])
			default:

			}
		}
	}()
}

func (b *Backend) SendMessage(text string) {
	b.conn.WriteToUDP([]byte(text), b.target)
}

func (b *Backend) Close() {
	b.conn.WriteToUDP([]byte("P2PLEAVMSG/"+b.username+"/"+b.ipHash), b.target)
	b.conn.Close()
}

func getIPHash() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	addr := conn.LocalAddr()
	if addr == nil {
		return "unknown"
	}

	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return "unknown"
	}

	return fmt.Sprintf("%x", md5.Sum(udpAddr.IP))
}
