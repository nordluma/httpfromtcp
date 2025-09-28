package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/nordluma/httpfromtcp/internal/request"
	"github.com/nordluma/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	status  response.StatusCode
	message string
}

func NewHandlerError(status response.StatusCode, message string) *HandlerError {
	return &HandlerError{
		status:  status,
		message: message,
	}
}

func (h HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.status)
	headers := response.GetDefaultHeaders(len(h.message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(h.message))
}

type Server struct {
	listener net.Listener
	handler  Handler

	closed atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
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
		handlerErr := NewHandlerError(response.BadRequest, err.Error())
		handlerErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	resErr := s.handler(buf, req)
	if resErr != nil {
		resErr.Write(conn)
		return
	}

	response.WriteStatusLine(conn, response.Ok)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, headers)
	conn.Write(buf.Bytes())
}
