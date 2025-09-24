package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

type requestState int

const (
	initialized requestState = iota
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	state       requestState
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		n, reqLine, err := parseRequestLine(data)
		if err != nil {
			// something bad happened
			return 0, err
		}

		if n == 0 {
			// need more data
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.state = done

		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIdx := 0
	req := &Request{state: initialized}

	for req.state != done {
		if readToIdx == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIdx:])
		if err != nil {
			if err == io.EOF {
				req.state = done
				break
			}

			return nil, err
		}
		readToIdx += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIdx])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIdx -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, nil, nil
	}

	reqLineStr := string(data[:idx])
	reqLine, err := requestLineFromString(reqLineStr)
	if err != nil {
		return 0, nil, err
	}

	return idx + 2, reqLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Malformed request-line: %s", str)
	}

	method, target, versionPart := parts[0], parts[1], parts[2]
	method, err := parseHttpMethod(method)
	if err != nil {
		return nil, err
	}

	version, err := parseHttpVersion(versionPart)
	if err != nil {
		return nil, err
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}

func parseHttpMethod(methodPart string) (string, error) {
	methodStr := strings.TrimSpace(methodPart)
	for _, char := range methodStr {
		if char < 'A' || char > 'Z' {
			return "", fmt.Errorf("Invalid method: %s", methodStr)
		}
	}

	return methodStr, nil
}

func parseHttpVersion(httpVersionPart string) (string, error) {
	parts := strings.Split(httpVersionPart, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("Malformed start-line: %s", httpVersionPart)
	}

	httpPart, versionPart := parts[0], parts[1]
	if httpPart != "HTTP" {
		return "", fmt.Errorf("Unrecognized HTTP-version: %s", httpPart)
	}

	if versionPart != "1.1" {
		return "", fmt.Errorf("Invalid HTTP version: %s", versionPart)
	}

	return versionPart, nil
}
