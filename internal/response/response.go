package response

import (
	"fmt"
	"io"

	"github.com/nordluma/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok StatusCode = iota
	BadRequest
	InternalError
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string
	switch statusCode {
	case Ok:
		statusLine = "HTTP/1.1 200 OK"
	case BadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case InternalError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	default:
		statusLine = ""
	}

	statusLine = fmt.Sprintf("%s\r\n", statusLine)
	if _, err := w.Write([]byte(statusLine)); err != nil {
		return err
	}

	return nil
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
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}

	return nil
}
