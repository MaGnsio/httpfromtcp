package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal("File not found, ", err)
	}

	log.Printf("Listening on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		log.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1) //channel cap means you can send this amount of messages without a receiver waiting
	/**
			If we remove go:
	        - we make the whole main function as one goroutine.
			- then func will wait for someone to receive messages from it.
			- and in main listen := getLinesChannel(f) will keep waiting for chanel to return.
			- which will cause a deadlock.

			But if we buffer the out chan to take all lines that we need to send. It will not block on send
	         and the rest of the code goes on.
		**/
	go func() {
		defer f.Close()
		defer close(out)
		var current_line string = ""
		for {
			const B = 8
			var data []byte = make([]byte, B)
			n, err := f.Read(data)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal("Something went wrong reading the file, ", err)
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				current_line += string(data[:i])
				out <- current_line
				current_line = string(data[i+1:])
			} else {
				current_line += string(data)
			}
		}
		if len(current_line) != 0 {
			out <- current_line
		}
	}()
	return out
}
