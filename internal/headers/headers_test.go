package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok := headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, 25, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:      localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok = headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, 35, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nXas: jax\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, ok = headers.Get("Host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", host)
	xas, ok := headers.Get("Xas")
	assert.True(t, ok)
	assert.Equal(t, "jax", xas)
	zxz, ok := headers.Get("zxz")
	assert.False(t, ok)
	assert.Equal(t, "", zxz)
	assert.Equal(t, 35, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Valid headers with multiple values
	headers = NewHeaders()
	data = []byte("Person: Tom\r\nPerson: Jax\r\nPerson: Lucifer\r\nPerson: Boby\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	person, ok := headers.Get("Person")
	assert.True(t, ok)
	assert.Equal(t, "Tom, Jax, Lucifer, Boby", person)
	assert.Equal(t, 59, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
