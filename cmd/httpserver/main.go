package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nordluma/httpfromtcp/internal/request"
	"github.com/nordluma/httpfromtcp/internal/response"
	"github.com/nordluma/httpfromtcp/internal/server"
)

const port = 42069

func defaultHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return server.NewHandlerError(
			response.BadRequest,
			"Your problem is not my problem\n",
		)
	case "/myproblem":
		return server.NewHandlerError(
			response.InternalError,
			"Woopsie, my bad\n",
		)
	}

	w.Write([]byte("All good, frfr\n"))
	return nil
}

func main() {
	server, err := server.Serve(port, defaultHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
