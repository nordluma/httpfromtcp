package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/nordluma/httpfromtcp/internal/request"
)

type Server struct {
	listener net.Listener
	closed   *atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}

			fmt.Printf("error accepting connection: %s\n", err.Error())
			continue
		}
		fmt.Printf("Connection accepted: %s\n", conn.RemoteAddr())

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("error parsing request: %s\n", err.Error())
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

	response := fmt.Sprintf("%s\r\n%s\r\n%s\r\n\r\n%s\n",
		"HTTP/1.1 200 OK",
		"Content-Type: text/plain",
		"Content-Length: 13",
		"Hello World!",
	)

	if _, err = conn.Write([]byte(response)); err != nil {
		fmt.Printf("error writing response: %s\n", err.Error())
	}
}
