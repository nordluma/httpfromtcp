package main

import (
	"fmt"
	"net"

	"github.com/nordluma/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Printf("error starting tcp listener: %s\n", err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("error accepting new connection: %s\n", err.Error())
			continue
		}
		fmt.Printf("Connection accepted: %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error parsing request: %s", err.Error())
			continue
		}

		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)
		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}
		fmt.Println("Body:")
		fmt.Printf("%s\n", string(req.Body))
	}
}
