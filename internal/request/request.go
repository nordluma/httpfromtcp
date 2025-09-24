package request

import (
	"errors"
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
	var request Request
	bytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("error reading request from reader: %s\n", err.Error())
		return nil, err
	}

	reqStr := string(bytes)
	requestLine, err := parseRequestLine(reqStr)
	if err != nil {
		fmt.Printf("error parsing request line part: %s\n", err.Error())
		return nil, err
	}

	request.RequestLine = requestLine

	return &request, nil
}

func parseRequestLine(reqStr string) (RequestLine, error) {
	var reqLine RequestLine
	parts := strings.Split(reqStr, "\r\n")
	reqLineParts := strings.Split(parts[0], " ")
	if len(reqLineParts) != 3 {
		return RequestLine{}, errors.New("Malformed request line part")
	}

	httpMethod, err := parseHttpMethod(reqLineParts[0])
	if err != nil {
		return RequestLine{}, err
	}

	reqLine.Method = httpMethod
	reqLine.RequestTarget = strings.TrimSpace(reqLineParts[1])

	httpVersion, err := parseHttpVersion(reqLineParts[2])
	if err != nil {
		return RequestLine{}, err
	}
	reqLine.HttpVersion = httpVersion

	return reqLine, nil
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
		return "", errors.New("malformed http version part")
	}

	if parts[1] != "1.1" {
		return "", fmt.Errorf("Invalid HTTP version: %s", parts[1])
	}

	return parts[1], nil
}
