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
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:      localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 33, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nXas: jax\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "jax", headers.Get("Xas"))
	assert.Equal(t, "", headers.Get("zxz"))
	assert.Equal(t, 33, n)
	assert.True(t, done, "Expected done to be false, got: %v", done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
