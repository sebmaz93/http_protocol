package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Failed to resolve udp address:", err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatal("Failed to dial UDP:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println("Error writing to UDP:", err)
			continue
		}
	}
}
