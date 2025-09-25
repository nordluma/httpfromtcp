package headers

import (
	"bytes"
	"fmt"
	"slices"
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
	key, err = parseHeaderKey(key)
	if err != nil {
		return 0, false, err
	}

	h.Set(key, strings.TrimSpace(value))

	// amount of bytes read is index + CRLF (2)
	return idx + 2, false, err
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	value, found := h[key]

	return value, found
}

var allowedSpecialChars = []rune{
	'!',
	'#',
	'$',
	'%',
	'&',
	'\'',
	'*',
	'+',
	'-',
	'.',
	'^',
	'_',
	'`',
	'|',
	'~',
}

func parseHeaderKey(key string) (string, error) {
	if key[len(key)-1] == ' ' {
		return "", fmt.Errorf("Invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)
	for _, char := range key {
		isUpperChar := char >= 'A' && char <= 'Z'
		isLowerChar := char >= 'a' && char <= 'z'
		isDigit := char >= '0' && char <= '9'
		isAllowedSpecialChar := slices.Contains(allowedSpecialChars, char)

		if isUpperChar || isLowerChar || isDigit || isAllowedSpecialChar {
			continue
		}

		return "", fmt.Errorf("Invalid header name: %s", key)
	}

	return key, nil
}
