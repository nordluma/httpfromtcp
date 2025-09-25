package request

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsesRequestLine(t *testing.T) {
	data := createRequest("GET / HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: len(data),
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestLineWithPath(t *testing.T) {
	data := createRequest("GET /coffee HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: len(data) / 2,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParsePostRequestLineWithPath(t *testing.T) {
	data := createRequest("POST /coffee HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 8,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestInvalidNumberOfRequestLineParts(t *testing.T) {
	data := createRequest("/coffee HTTP/1.1")
	_, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 5,
	},
	)
	require.Error(t, err)
}

func TestInvalidMethodOrder(t *testing.T) {
	data := createRequest("/ GET HTTP/1.1")
	_, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 50,
	},
	)
	require.Error(t, err)
}

func TestInvalidHttpVersion(t *testing.T) {
	data := createRequest("GET / HTTP/69.0")
	_, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 10,
	})
	require.Error(t, err)
}

func TestReadRequestWithThreeByteChunks(t *testing.T) {
	data := createRequest("GET / HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 3,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestReadRequestWithOneByteChunks(t *testing.T) {
	data := createRequest("GET /coffee HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestParseRequestAndStandardHeaders(t *testing.T) {
	data := createRequest("GET / HTTP/1.1")
	r, err := RequestFromReader(&chunkReader{
		data:            data,
		numBytesPerRead: 3,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])
}

func TestParseRequestWithEmptyHeaders(t *testing.T) {
	r, err := RequestFromReader(&chunkReader{
		data:            "GET /coffee HTTP/1.1\r\n\r\n",
		numBytesPerRead: 5,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
}

func TestParseRequestWithMalformedHeader(t *testing.T) {
	_, err := RequestFromReader(&chunkReader{
		data: fmt.Sprintf(
			"%s\r\n%s\r\n\r\n",
			"GET / HTTP/1.1",
			"Host localhost:42069",
		),
		numBytesPerRead: 3,
	})
	assert.Error(t, err)
}

func TestParseRequestWithDuplicateHeaders(t *testing.T) {
	r, err := RequestFromReader(&chunkReader{
		data: fmt.Sprintf(
			"%s\r\n%s\r\n%s\r\n\r\n",
			"GET / HTTP/1.1",
			"Accept-Encoding: gzip",
			"Accept-Encoding: brotli",
		),
		numBytesPerRead: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "gzip, brotli", r.Headers["accept-encoding"])
}

func TestParseRequestWithCaseInsensitiveHeaders(t *testing.T) {
	r, err := RequestFromReader(&chunkReader{
		data: fmt.Sprintf(
			"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
			"GET / HTTP/1.1",
			"HOST: localhost:42069",
			"USER-AGENT: curl/7.81",
			"ACCEPT: */*",
		),
		numBytesPerRead: 15,
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])
}

func TestParseRequestMissingEOFHeaders(t *testing.T) {
	_, err := RequestFromReader(&chunkReader{
		data: fmt.Sprintf(
			"%s\r\n%s",
			"POST /password HTTP/1.1",
			"Host: localhost:42069",
		),
		numBytesPerRead: 1,
	})
	require.Error(t, err)
}

func createRequest(requestLine string) string {
	return fmt.Sprintf(
		"%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		requestLine,
		"Host: localhost:42069",
		"User-Agent: curl/7.81",
		"Accept: */*",
	)
}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIdx := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, []byte(cr.data)[cr.pos:endIdx])
	cr.pos += n

	return n, nil
}
