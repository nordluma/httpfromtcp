package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(bytes)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, fmt.Errorf("Failed to find CRLF in request-line")
	}

	reqLineStr := string(data[:idx])
	reqLine, err := requestLineFromString(reqLineStr)
	if err != nil {
		return nil, err
	}

	return reqLine, nil
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
		if !unicode.IsLetter(char) || !unicode.IsUpper(char) {
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
