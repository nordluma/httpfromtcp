package response

import (
	"fmt"

	"github.com/nordluma/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) string {
	reasonPhrase := ""
	switch statusCode {
	case Ok:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalError:
		reasonPhrase = "Internal Server Error"
	}

	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}
