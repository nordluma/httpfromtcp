package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nordluma/httpfromtcp/internal/request"
	"github.com/nordluma/httpfromtcp/internal/response"
	"github.com/nordluma/httpfromtcp/internal/server"
)

const port = 42069

func defaultHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		handler400(w, req)
		return
	}

	if target == "/myproblem" {
		handler500(w, req)
		return
	}

	if strings.HasPrefix(target, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	handler200(w, req)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	res, err := http.Get("https://httpbin.org/" + target)
	if err != nil {
		handler500(w, req)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.Ok)
	h := response.GetDefaultHeaders(0)
	h.Replace("transfer-encoding", "chunked")
	h.Delete("content-length")
	w.WriteHeaders(h)

	const chunkSize = 1024
	for {
		buf := make([]byte, chunkSize)
		n, err := res.Body.Read(buf)
		fmt.Printf("read %d bytes from stream\n", n)
		if n > 0 {
			_, err = w.WriteChunkedBody(buf)
			if err != nil {
				fmt.Printf("error writing chunk: %v", err)
				break
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("error reading response body: %v\n", err)
			break
		}
	}

	if _, err = w.WriteChunkedBodyDone(); err != nil {
		fmt.Printf("error writing chunked body done: %v\n", err)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BadRequest)
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)

	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.InternalError)
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)

	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.Ok)
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)

	h := response.GetDefaultHeaders(len(body))
	h.Replace("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
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
