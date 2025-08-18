package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("failed to resolve UDP address: ", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("failed to establish UDP connection: ", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("failed to read from stdin: ", err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatal("failed to write to connection: ", err)
		}
	}
}
