package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nordluma/httpfromtcp/internal/headers"
	"github.com/nordluma/httpfromtcp/internal/request"
	"github.com/nordluma/httpfromtcp/internal/response"
	"github.com/nordluma/httpfromtcp/internal/server"
)

const port = 42069

func defaultHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	if target == "/video" {
		videoHandler(w, req)
		return
	}

	if target == "/yourproblem" {
		handler400(w, req)
		return
	}

	if target == "/myproblem" {
		handler500(w, req)
		return
	}

	handler200(w, req)
}

func videoHandler(w *response.Writer, req *request.Request) {
	videoBytes, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		fmt.Printf("error reading video: %v", err)
		handler500(w, req)
		return
	}

	w.WriteStatusLine(response.Ok)
	h := response.GetDefaultHeaders(len(videoBytes))
	h.Override("content-type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(videoBytes)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Printf("Proxying to: %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.Ok)
	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-SHA256, X-Content-Length")
	h.Delete("content-length")
	w.WriteHeaders(h)

	const chunkSize = 1024
	fullBody := make([]byte, 0)
	buf := make([]byte, chunkSize)
	for {
		n, err := res.Body.Read(buf)
		fmt.Printf("read %d bytes\n", n)
		if n > 0 {
			if _, err = w.WriteChunkedBody(buf[:n]); err != nil {
				fmt.Printf("error writing chunk: %v\n", err)
				break
			}

			fullBody = append(fullBody, buf[:n]...)
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

	trailers := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailers.Override("X-Content-SHA256", sha256)
	trailers.Override("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

	if err = w.WriteTrailers(trailers); err != nil {
		fmt.Printf("error writing trailers: %v", err)
	}

	fmt.Println("trailers written")
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
	h.Override("Content-Type", "text/html")
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
	h.Override("Content-Type", "text/html")
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
	h.Override("Content-Type", "text/html")
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
