package main

import (
	"fmt"
	"log"
	"net"
	"tcpToHttp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	defer listener.Close()

	if err != nil {
		log.Fatal("error", err)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}
		fmt.Println("connection has been accepted.")
		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatal("error", err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		req.Headers.ForEach(func(k, v string) {
			fmt.Printf("- %s: %s\n", k, v)
		})
		fmt.Printf("Body:\n")
		fmt.Printf(string(req.Body))
	}

}
