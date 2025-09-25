package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, found := headers.Get("Host")
	assert.True(t, found)
	assert.Equal(t, "localhost:42069", value)
	assert.Equal(t, len(data)-2, n) // the last CRLF is excluded
	assert.False(t, done)
}

func TestInvalidSpacing(t *testing.T) {
	headers := NewHeaders()
	data := []byte("    Host : localhost:42069        \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseSingleHeaderWithExtraSpace(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069 \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, found := headers.Get("Host")
	assert.True(t, found)
	assert.Equal(t, "localhost:42069", value)
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)
}

func TestParseTwoHeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Encoding", "gzip")
	data := []byte("Host: localhost:42069\r\nContent-Length: 55\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, found := headers.Get("Host")
	assert.True(t, found)
	assert.Equal(t, "localhost:42069", value)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	data = data[n:] // move to the next header
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, found = headers.Get("Content-Length")
	assert.True(t, found)
	assert.Equal(t, "55", value)
	assert.Equal(t, 20, n)
	assert.False(t, done)

	value, found = headers.Get("Content-Type")
	assert.True(t, found)
	assert.Equal(t, "application/json", value)

	value, found = headers.Get("Accept-Encoding")
	assert.True(t, found)
	assert.Equal(t, "gzip", value)
}

func TestParseDoneFieldLine(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestInvalidCharacterInHeaderKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
