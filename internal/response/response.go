package response

import (
	"fmt"
	"io"

	"github.com/nordluma/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase := ""
	switch statusCode {
	case Ok:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalError:
		reasonPhrase = "Internal Server Error"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)

	_, err := w.Write([]byte(statusLine))

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		header := fmt.Sprintf("%s: %s\r\n", key, val)
		if _, err := w.Write([]byte(header)); err != nil {
			return err
		}
	}

	// add the end of headers
	_, err := w.Write([]byte("\r\n"))

	return err
}
