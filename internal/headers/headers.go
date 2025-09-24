package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		// need more data
		return 0, false, nil
	}

	if idx == 0 {
		// no headers found, we're done -> consume the CRLF
		return 2, true, nil
	}

	headerStr := string(data[:idx])
	pair := strings.SplitN(headerStr, ":", 2)

	key, value := pair[0], pair[1]
	if key[len(key)-1] == ' ' {
		return 0, false, fmt.Errorf("Invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)
	h[key] = strings.TrimSpace(value)

	// amount of bytes read is index + CRLF (2)
	return idx + 2, false, err
}
