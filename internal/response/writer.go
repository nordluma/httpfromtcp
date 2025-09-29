package response

import (
	"fmt"
	"io"

	"github.com/nordluma/httpfromtcp/internal/headers"
)

type writerState int

const (
	stateStatusLine writerState = iota
	stateHeaders
	stateBody
	stateTrailers
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  stateStatusLine,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateStatusLine {
		return fmt.Errorf("cannot write status line in state: %d", w.state)
	}
	defer func() { w.state = stateHeaders }()

	statusLine := getStatusLine(statusCode)
	_, err := w.writer.Write([]byte(statusLine))

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateHeaders {
		return fmt.Errorf("cannot write headers in state: %d", w.state)
	}
	defer func() { w.state = stateBody }()

	for key, val := range headers {
		header := fmt.Sprintf("%s: %s\r\n", key, val)
		if _, err := w.writer.Write([]byte(header)); err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.state)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.state)
	}

	chunkSize := len(p)
	total := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return total, err
	}
	total += n

	n, err = w.writer.Write(p)
	if err != nil {
		return total, err
	}
	total += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return total, err
	}
	total += n

	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.state)
	}
	defer func() { w.state = stateTrailers }()

	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}
