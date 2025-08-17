package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer file.Close()
		defer close(ch)

		var line []byte
		for {
			data := make([]byte, 8)
			_, err := file.Read(data)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("failed to read data: ", err)
			}

			parts := bytes.Split(data, []byte("\n"))
			for idx := range parts {
				line = append(line, parts[idx]...)
				if len(parts[idx:]) > 1 {
					ch <- string(line)
					line = nil
				}
			}
		}
	}()

	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("failed to open file: ", err)
	}

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Println("read:", line)
	}
}
