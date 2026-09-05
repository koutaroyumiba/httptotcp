package main

import (
	"fmt"
	"log"
	"net"

	"github.com/koutaroyumiba/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Failed to listen to network: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to establish connection: %v", err)
		}
		fmt.Println("Established connection")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Failed to read request: %v", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		for k, v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Closing connection...")
		conn.Close()
	}
}
