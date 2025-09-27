package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/nordluma/httpfromtcp/internal/request"
	"github.com/nordluma/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
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

	if err = response.WriteStatusLine(conn, response.Ok); err != nil {
		fmt.Printf("error writing status line: %s\n", err.Error())
	}

	headers := response.GetDefaultHeaders(0)
	if err = response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error writing headers: %s\n", err)
	}
}
