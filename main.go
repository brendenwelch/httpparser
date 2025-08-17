package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("failed to open file:", err)
	}

	var line []byte
	for {
		data := make([]byte, 8)
		_, err := file.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("failed to read data:", err)
		}

		parts := bytes.Split(data, []byte("\n"))
		for idx := range parts {
			line = append(line, parts[idx]...)
			if len(parts[idx:]) > 1 {
				fmt.Println("read:", string(line))
				line = nil
			}
		}
	}
}
