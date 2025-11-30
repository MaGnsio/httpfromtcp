package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	var serverAddr string = "localhost:42069"
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		message, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("message sent\n")
	}
}
