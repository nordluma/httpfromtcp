package request

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsesRequestLine(t *testing.T) {
	r, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"GET / HTTP/1.1",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLineWithPath(t *testing.T) {
	r, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"GET /coffee HTTP/1.1",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParsePostRequestLineWithPath(t *testing.T) {
	r, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"POST /coffee HTTP/1.1",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestInvalidNumberOfRequestLineParts(t *testing.T) {
	_, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"/coffee HTTP/1.1",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.Error(t, err)
}

func TestInvalidMethodOrder(t *testing.T) {
	_, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"/ GET HTTP/1.1",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.Error(t, err)
}

func TestInvalidHttpVersion(t *testing.T) {
	_, err := RequestFromReader(
		strings.NewReader(
			fmt.Sprintf(
				"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
				"GET / HTTP/69.0",
				"Host: localhost:42069",
				"User-Agent: curl/7.81",
				"Accept: */*",
			),
		),
	)
	require.Error(t, err)
}
