package main

import (
	"fmt"
	"log"
	"net"

	"github.com/brendenwelch/httpparser/internal/request"
)

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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("failed to process request from reader: ", err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)
	}
}
