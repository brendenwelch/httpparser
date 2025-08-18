package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(stream io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer stream.Close()
		defer close(out)

		var line []byte
		for {
			data := make([]byte, 8)
			_, err := stream.Read(data)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("failed to read data: ", err)
			}

			parts := bytes.Split(data, []byte("\n"))
			for idx := range parts {
				line = append(line, parts[idx]...)
				if len(parts[idx:]) > 1 {
					out <- string(line)
					line = nil
				}
			}
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatal("failed to open tcp listener: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept connection from listener: ", err)
		}

		fmt.Println("Connection accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
