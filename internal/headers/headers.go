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
		// no headers found, we're done
		return 0, true, nil
	}

	headerStr := string(data[:idx])
	n = len(headerStr) + 2
	headerPair := strings.SplitN(headerStr, ":", 2)
	if len(headerPair) != 2 {
		return 0, false, fmt.Errorf(
			"Malformed header pair: %v, len: %d",
			headerPair,
			len(headerPair),
		)
	}

	key, value := headerPair[0], strings.TrimSpace(headerPair[1])
	if key[len(key)-1] == ' ' {
		return 0, false, fmt.Errorf("Malformed header key: %s", key)
	}

	key = strings.TrimSpace(key)
	h[key] = value

	return n, done, err
}
